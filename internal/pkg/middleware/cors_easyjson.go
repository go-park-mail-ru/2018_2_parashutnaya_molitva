// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package middleware

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonE7f5bc6fDecodeGithubComGoParkMailRu20182ParashutnayaMolitvaInternalPkgMiddleware(in *jlexer.Lexer, out *CorsData) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "AllowOrigins":
			if in.IsNull() {
				in.Skip()
				out.AllowOrigins = nil
			} else {
				in.Delim('[')
				if out.AllowOrigins == nil {
					if !in.IsDelim(']') {
						out.AllowOrigins = make([]string, 0, 4)
					} else {
						out.AllowOrigins = []string{}
					}
				} else {
					out.AllowOrigins = (out.AllowOrigins)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.AllowOrigins = append(out.AllowOrigins, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "AllowMethods":
			if in.IsNull() {
				in.Skip()
				out.AllowMethods = nil
			} else {
				in.Delim('[')
				if out.AllowMethods == nil {
					if !in.IsDelim(']') {
						out.AllowMethods = make([]string, 0, 4)
					} else {
						out.AllowMethods = []string{}
					}
				} else {
					out.AllowMethods = (out.AllowMethods)[:0]
				}
				for !in.IsDelim(']') {
					var v2 string
					v2 = string(in.String())
					out.AllowMethods = append(out.AllowMethods, v2)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "AllowHeaders":
			if in.IsNull() {
				in.Skip()
				out.AllowHeaders = nil
			} else {
				in.Delim('[')
				if out.AllowHeaders == nil {
					if !in.IsDelim(']') {
						out.AllowHeaders = make([]string, 0, 4)
					} else {
						out.AllowHeaders = []string{}
					}
				} else {
					out.AllowHeaders = (out.AllowHeaders)[:0]
				}
				for !in.IsDelim(']') {
					var v3 string
					v3 = string(in.String())
					out.AllowHeaders = append(out.AllowHeaders, v3)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "MaxAge":
			out.MaxAge = int(in.Int())
		case "AllowCredentials":
			out.AllowCredentials = bool(in.Bool())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonE7f5bc6fEncodeGithubComGoParkMailRu20182ParashutnayaMolitvaInternalPkgMiddleware(out *jwriter.Writer, in CorsData) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"AllowOrigins\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.AllowOrigins == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v4, v5 := range in.AllowOrigins {
				if v4 > 0 {
					out.RawByte(',')
				}
				out.String(string(v5))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"AllowMethods\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.AllowMethods == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v6, v7 := range in.AllowMethods {
				if v6 > 0 {
					out.RawByte(',')
				}
				out.String(string(v7))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"AllowHeaders\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.AllowHeaders == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.AllowHeaders {
				if v8 > 0 {
					out.RawByte(',')
				}
				out.String(string(v9))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"MaxAge\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.MaxAge))
	}
	{
		const prefix string = ",\"AllowCredentials\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.AllowCredentials))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CorsData) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonE7f5bc6fEncodeGithubComGoParkMailRu20182ParashutnayaMolitvaInternalPkgMiddleware(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CorsData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonE7f5bc6fEncodeGithubComGoParkMailRu20182ParashutnayaMolitvaInternalPkgMiddleware(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CorsData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonE7f5bc6fDecodeGithubComGoParkMailRu20182ParashutnayaMolitvaInternalPkgMiddleware(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CorsData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonE7f5bc6fDecodeGithubComGoParkMailRu20182ParashutnayaMolitvaInternalPkgMiddleware(l, v)
}
