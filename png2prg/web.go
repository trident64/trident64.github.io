package main

import (
	"bytes"
	"net/http"
	"syscall/js"

	"github.com/staD020/png2prg"
)

func main() {
	// Set up the global JS function
	js.Global().Set("convertPngToPrg", js.FuncOf(convertPngToPrg))

	// Keep the Go program running
	<-make(chan bool)
}

// convertPngToPrg is exposed to JavaScript
func convertPngToPrg(this js.Value, args []js.Value) interface{} {
	// Check if we have the correct number of arguments
	if len(args) < 1 {
		return map[string]interface{}{
			"error": "Missing URL argument",
		}
	}

	// Get URL from JavaScript
	url := args[0].String()

	// Parse options from JavaScript if provided
	options := png2prg.Options{}
	if len(args) > 1 && !args[1].IsNull() && !args[1].IsUndefined() {
		jsOptions := args[1]

		// Example of setting some common options
		if jsOptions.Get("mode").Type() == js.TypeString {
			options.Mode = jsOptions.Get("mode").String()
		}

		if jsOptions.Get("bitpairColors").Type() == js.TypeString {
			options.BitpairColors = jsOptions.Get("bitpairColors").String()
		}

		if jsOptions.Get("display").Type() == js.TypeBoolean {
			options.Display = jsOptions.Get("display").Bool()
		}

		// Add more options as needed
	}

	// Fetch the image
	resp, err := http.Get(url)
	if err != nil {
		return map[string]interface{}{
			"error": "Failed to fetch image: " + err.Error(),
		}
	}
	defer resp.Body.Close()

	// Create PNG2PRG processor
	p, err := png2prg.New(options, resp.Body)
	if err != nil {
		return map[string]interface{}{
			"error": "Failed to process image: " + err.Error(),
		}
	}

	// Convert to PRG
	var buf bytes.Buffer
	_, err = p.WriteTo(&buf)
	if err != nil {
		return map[string]interface{}{
			"error": "Failed to convert image: " + err.Error(),
		}
	}

	// Convert the buffer to a JavaScript Uint8Array
	data := buf.Bytes()
	jsArray := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(jsArray, data)

	// Return the result
	return map[string]interface{}{
		"data": jsArray,
		"size": len(data),
	}
}
