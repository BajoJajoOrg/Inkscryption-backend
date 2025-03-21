// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package images

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

func easyjsonEafed2a8DecodeGithubComBajoJajoOrgInkscryptionBackendImages(in *jlexer.Lexer, out *Image) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "person_id":
			out.UserId = int64(in.Int64())
		case "image_url":
			out.Url = string(in.String())
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
func easyjsonEafed2a8EncodeGithubComBajoJajoOrgInkscryptionBackendImages(out *jwriter.Writer, in Image) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"person_id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.UserId))
	}
	{
		const prefix string = ",\"image_url\":"
		out.RawString(prefix)
		out.String(string(in.Url))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Image) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonEafed2a8EncodeGithubComBajoJajoOrgInkscryptionBackendImages(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Image) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonEafed2a8EncodeGithubComBajoJajoOrgInkscryptionBackendImages(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Image) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonEafed2a8DecodeGithubComBajoJajoOrgInkscryptionBackendImages(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Image) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonEafed2a8DecodeGithubComBajoJajoOrgInkscryptionBackendImages(l, v)
}
func easyjsonEafed2a8DecodeGithubComBajoJajoOrgInkscryptionBackendImages1(in *jlexer.Lexer, out *Canvases) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "canvases":
			if in.IsNull() {
				in.Skip()
				out.Canvases = nil
			} else {
				in.Delim('[')
				if out.Canvases == nil {
					if !in.IsDelim(']') {
						out.Canvases = make([]Canvas, 0, 1)
					} else {
						out.Canvases = []Canvas{}
					}
				} else {
					out.Canvases = (out.Canvases)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Canvas
					(v1).UnmarshalEasyJSON(in)
					out.Canvases = append(out.Canvases, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjsonEafed2a8EncodeGithubComBajoJajoOrgInkscryptionBackendImages1(out *jwriter.Writer, in Canvases) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"canvases\":"
		out.RawString(prefix[1:])
		if in.Canvases == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Canvases {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Canvases) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonEafed2a8EncodeGithubComBajoJajoOrgInkscryptionBackendImages1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Canvases) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonEafed2a8EncodeGithubComBajoJajoOrgInkscryptionBackendImages1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Canvases) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonEafed2a8DecodeGithubComBajoJajoOrgInkscryptionBackendImages1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Canvases) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonEafed2a8DecodeGithubComBajoJajoOrgInkscryptionBackendImages1(l, v)
}
func easyjsonEafed2a8DecodeGithubComBajoJajoOrgInkscryptionBackendImages2(in *jlexer.Lexer, out *Canvas) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "canvas_name":
			out.Name = string(in.String())
		case "canvas_url":
			out.Url = string(in.String())
		case "update_time":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Update).UnmarshalJSON(data))
			}
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
func easyjsonEafed2a8EncodeGithubComBajoJajoOrgInkscryptionBackendImages2(out *jwriter.Writer, in Canvas) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"canvas_name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"canvas_url\":"
		out.RawString(prefix)
		out.String(string(in.Url))
	}
	{
		const prefix string = ",\"update_time\":"
		out.RawString(prefix)
		out.Raw((in.Update).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Canvas) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonEafed2a8EncodeGithubComBajoJajoOrgInkscryptionBackendImages2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Canvas) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonEafed2a8EncodeGithubComBajoJajoOrgInkscryptionBackendImages2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Canvas) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonEafed2a8DecodeGithubComBajoJajoOrgInkscryptionBackendImages2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Canvas) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonEafed2a8DecodeGithubComBajoJajoOrgInkscryptionBackendImages2(l, v)
}
