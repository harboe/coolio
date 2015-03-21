package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const minifyURL = "https://closure-compiler.appspot.com/compile"

// {
// "compiledCode":/* raw code here */,
// {"errors": [
//   {"charno":4321,
//    "error":"ERROR: You failed.",
//    "lineno":1234,
//    "file":"default.js",
//    "type":"ERROR_TYPE",
//    "line":"var x=-'hello';"}],
// "warnings": [
//   {"charno":4321,
//    "lineno":1234,
//    "file":"default.js",
//    "type":"ERROR_TYPE",
//    "warning":"Warning: You did something wrong!",
//    "line":"delete 1;"}]
// "serverErrors":[
//   {"code":123,"error":"Over quota"}
//   ],
// "statistics":{
//   "originalSize":10,
//   "compressedSize":3000
//   "compileTime":10
//   }
// }
type closureResponse struct {
}

func Minify(jsCode []byte) ([]byte, error) {
	form := url.Values{
		"js_code":           []string{string(jsCode)},
		"compilation_level": []string{"SIMPLE_OPTIMIZATIONS"},
		"output_format":     []string{"json"},
		"output_info":       []string{"compiled_code", "errors", "statistics"},
	}

	resp, err := http.PostForm(minifyURL, form)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	var out map[string]interface{}
	err = decoder.Decode(&out)
	if err != nil {
		return []byte{}, err
	}
	if se := out["serverErrors"]; se != nil {
		return []byte{}, fmt.Errorf("server errors: %v", se)
	}
	if e := out["errors"]; e != nil {
		return []byte{}, fmt.Errorf("errors: %v", e)
	}
	if w := out["warning"]; w != nil {
		log.Println("warnings when compiling code:", w)
	}
	if stats := out["statistics"]; stats != nil {
		if statsm, _ := stats.(map[string]interface{}); statsm != nil {
			log.Printf("Compressed JS from %v to %v (GZIP'ed %v to %v)",
				statsm["originalSize"], statsm["compressedSize"],
				statsm["originalGzipSize"], statsm["compressedGzipSize"])
		}
	}

	// log.Println(out["compiledCode"])
	// compiled := out["compiledCode"].(string)
	// _, err = io.WriteString(w, compiled)
	return []byte(out["compiledCode"].(string)), err
}
