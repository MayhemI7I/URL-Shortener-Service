package zstd

import (
	"io"
	"local/logger"
	"net/http"
	"strings"

	"github.com/klauspost/compress/zstd"
)



type ZstdRspWrt struct {
	http.ResponseWriter
	Writer *zstd.Encoder
}

func (zsw *ZstdRspWrt) Write(b []byte)(int, error){
	return zsw.Writer.Write(b)
}

func ZstdDecompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "zstd") {
			next.ServeHTTP(w, r)
			return
		}

		decoder, _ := zstd.NewReader(r.Body)

		defer decoder.Close()

		r.Body = io.NopCloser(decoder)
		w.Header().Del("Content-Length")
		logger.Log.Debug("The response decompression procedure has been initialized")

		next.ServeHTTP(w, r)
	})
}

 

func ZstdCompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "zstd") {
			next.ServeHTTP(w, r)
			return
		}

		encoder, _ := zstd.NewWriter(w, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
		
		defer encoder.Close()
		w.Header().Set("Content-Encoding", "zstd")
		w.Header().Del("Content-Lenght")

		logger.Log.Debug("The response compression procedure has been initialized")

		next.ServeHTTP(&ZstdRspWrt{ResponseWriter: w, Writer: encoder}, r)

	})
}

