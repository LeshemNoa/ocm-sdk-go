/*
Copyright (c) 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

import (
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalAddOnRequirementStatus writes a value of the 'add_on_requirement_status' type to the given writer.
func MarshalAddOnRequirementStatus(object *AddOnRequirementStatus, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeAddOnRequirementStatus(object, stream)
	stream.Flush()
	return stream.Error
}

// writeAddOnRequirementStatus writes a value of the 'add_on_requirement_status' type to the given stream.
func writeAddOnRequirementStatus(object *AddOnRequirementStatus, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.errorMsgs != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("error_msgs")
		writeStringList(object.errorMsgs, stream)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("fulfilled")
		stream.WriteBool(object.fulfilled)
		count++
	}
	stream.WriteObjectEnd()
}

// UnmarshalAddOnRequirementStatus reads a value of the 'add_on_requirement_status' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAddOnRequirementStatus(source interface{}) (object *AddOnRequirementStatus, err error) {
	if source == http.NoBody {
		return
	}
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readAddOnRequirementStatus(iterator)
	err = iterator.Error
	return
}

// readAddOnRequirementStatus reads a value of the 'add_on_requirement_status' type from the given iterator.
func readAddOnRequirementStatus(iterator *jsoniter.Iterator) *AddOnRequirementStatus {
	object := &AddOnRequirementStatus{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "error_msgs":
			value := readStringList(iterator)
			object.errorMsgs = value
			object.bitmap_ |= 1
		case "fulfilled":
			value := iterator.ReadBool()
			object.fulfilled = value
			object.bitmap_ |= 2
		default:
			iterator.ReadAny()
		}
	}
	return object
}
