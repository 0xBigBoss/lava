package main

import (
	"fmt"
	"regexp"
	"strconv"
)

func events() map[string]func(LogLine) TestResult {
	tests := map[string](func(LogLine) TestResult){
		"🔄":                          test_start,
		"🌍":                          test_found_pass,
		"lava_spec_add":              test_found_pass,
		"lava_provider_stake_new":    test_found_pass,
		"lava_client_stake_new":      test_found_pass,
		"lava_relay_payment":         test_found_pass,
		"ERR_client_entries_pairing": test_ERR_client_entries_pairing,
		"update pairing list!":       test_found_pass,
		"Client pubkey":              test_found_pass,
		"no pairings available":      test_found_fail,
		"rpc error":                  test_found_pass,
		"reply":                      test_found_pass,
		"refused":                    test_found_fail,
		"listening":                  test_found_pass,
		"init done":                  test_found_pass,
		"connection refused":         test_found_fail_now,
		"cannot build app":           test_found_fail_now,
		"exit status":                test_found_fail_now,
	}
	return tests
}

func lava_up(line string, params []interface{}) TestResult {
	contains := "Token faucet"
	return test_basic(line, contains)
}
func init_done(line string, params []interface{}) TestResult {
	contains := "init done"
	return test_basic(line, contains)
}
func raw_log(line string, params []interface{}) TestResult {
	contains := "raw_log"
	return test_basic(line, contains)
}
func providers_ready(line string, params []interface{}) TestResult {
	contains := "listening"
	return test_basic(line, contains)
}

func providers_ready_eth(line string, params []interface{}) TestResult {
	contains := "starting"
	return test_basic(line, contains)
}

func found_rpc_reply(line string, params []interface{}) TestResult {
	contains := "reply JSONRPC_"
	return test_basic(line, contains)
}

func client_finished(line string, params []interface{}) TestResult {
	contains := "Client finished"
	return test_basic(line, contains)
}

func found_relay_payment(line string, params []interface{}) TestResult {
	contains := "lava_relay_payment"
	if len(params) != 9 {
		panic("found_relay_payment not enough params")
	}
	var testResult TestResult
	tr1 := test_float_param(line, contains, []interface{}{params[0], params[1], params[2]})
	tr2 := test_float_param(line, contains, []interface{}{params[3], params[4], params[5]})
	tr3 := test_float_param(line, contains, []interface{}{params[6], params[7], params[8]})

	testResult.eventID = contains
	errorstring := ""
	if tr1.err != nil {
		errorstring += tr1.err.Error()
	}
	if tr2.err != nil {
		errorstring += tr2.err.Error()
	}
	if tr3.err != nil {
		errorstring += tr3.err.Error()
	}
	testResult.err = fmt.Errorf(errorstring)
	if errorstring == "" {
		testResult.err = nil
	}

	testResult.failNow = tr1.failNow && tr2.failNow && tr3.failNow
	testResult.found = tr1.found && tr2.found && tr3.found
	testResult.passed = tr1.passed && tr2.passed && tr3.passed
	testResult.line = line
	return testResult
}
func osmosis_finished(line string, params []interface{}) TestResult {
	contains := "osmosis finished"
	return test_basic(line, contains)
}
func node_reset(line string, params []interface{}) TestResult {
	contains := "🔄"
	return test_basic(line, contains)
}
func node_ready(line string, params []interface{}) TestResult {
	contains := "🌍 Token faucet: http"
	return test_basic(line, contains)
}
func new_epoch(line string, params []interface{}) TestResult {
	contains := "lava_new_epoch"
	return test_basic(line, contains)
}

func test_found_pass(log LogLine) TestResult {
	return TestResult{
		eventID: "found_pass",
		found:   true,
		passed:  true,
		line:    log.line,
		err:     nil,
		parent:  log.parent,
		failNow: false,
	}
}

func test_found_fail(log LogLine) TestResult {
	return TestResult{
		eventID: "found_fail",
		found:   true,
		passed:  false,
		line:    log.line,
		err:     nil,
		parent:  log.parent,
		failNow: false,
	}
}

func test_found_fail_now(log LogLine) TestResult {
	return TestResult{
		eventID: "found_fail_now",
		found:   true,
		passed:  false,
		line:    log.line,
		err:     nil,
		parent:  log.parent,
		failNow: true,
	}
}

func test_basic(line string, contains string) TestResult {
	found, pass := false, false
	if strContains(line, contains) {
		found, pass = true, true
	}
	return TestResult{
		eventID: "",
		found:   found,
		passed:  pass,
		line:    line,
		err:     nil,
		parent:  "",
		failNow: false,
	}
}

func test_float_param(line string, contains string, params []interface{}) TestResult {
	found, pass := false, false
	var err error
	if len(params) == 3 {
		if strContains(line, contains) {
			found = true
			paramstring, okstr := params[0].(string)
			lower, oklower := params[1].(float64)
			higher, okhigher := params[2].(float64)
			if okstr && oklower && okhigher {
				re := regexp.MustCompile(paramstring + "(\\d.\\d+)")
				if re.MatchString(line) {
					res, okres := strconv.ParseFloat(re.FindStringSubmatch(line)[1], 64)
					if okres == nil {
						if lower <= res && res <= higher {
							pass = true
							err = nil
						} else {
							err = fmt.Errorf("%s: %f out of expected range [%f,%f]", paramstring, res, lower, higher)
						}
					}
				}
			}
		}
	}
	return TestResult{
		eventID: contains,
		found:   found,
		passed:  pass,
		line:    line,
		err:     err,
		parent:  "",
		failNow: false,
	}
}

func test_ERR_client_entries_pairing(log LogLine) TestResult {
	return TestResult{
		eventID: "found_fail_now",
		found:   true,
		passed:  true,
		line:    log.line,
		err:     fmt.Errorf("ERR_client_entries_pairing is unexpected but still passing to finish fullflow"),
		parent:  log.parent,
		failNow: false,
	}
}

func test_start(log LogLine) TestResult {
	return TestResult{
		eventID: "found_fail_now",
		found:   true,
		passed:  true,
		line:    log.line,
		err:     fmt.Errorf("🔄 is not expected"),
		parent:  log.parent,
		failNow: false,
	}
}

func test_start_fail(log LogLine) TestResult {
	return TestResult{
		eventID: "found_fail_now",
		found:   true,
		passed:  false,
		line:    log.line,
		err:     fmt.Errorf("🔄 is not expected"),
		parent:  log.parent,
		failNow: true,
	}
}
