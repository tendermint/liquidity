/* eslint-disable */
import { Writer, Reader } from "protobufjs/minimal";
export const protobufPackage = "google.protobuf";
/**
 * `NullValue` is a singleton enumeration to represent the null value for the
 * `Value` type union.
 *
 *  The JSON representation for `NullValue` is JSON `null`.
 */
export var NullValue;
(function (NullValue) {
    /** NULL_VALUE - Null value. */
    NullValue[NullValue["NULL_VALUE"] = 0] = "NULL_VALUE";
    NullValue[NullValue["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(NullValue || (NullValue = {}));
export function nullValueFromJSON(object) {
    switch (object) {
        case 0:
        case "NULL_VALUE":
            return NullValue.NULL_VALUE;
        case -1:
        case "UNRECOGNIZED":
        default:
            return NullValue.UNRECOGNIZED;
    }
}
export function nullValueToJSON(object) {
    switch (object) {
        case NullValue.NULL_VALUE:
            return "NULL_VALUE";
        default:
            return "UNKNOWN";
    }
}
const baseStruct = {};
export const Struct = {
    encode(message, writer = Writer.create()) {
        Object.entries(message.fields).forEach(([key, value]) => {
            Struct_FieldsEntry.encode({ key: key, value }, writer.uint32(10).fork()).ldelim();
        });
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseStruct };
        message.fields = {};
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    const entry1 = Struct_FieldsEntry.decode(reader, reader.uint32());
                    if (entry1.value !== undefined) {
                        message.fields[entry1.key] = entry1.value;
                    }
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseStruct };
        message.fields = {};
        if (object.fields !== undefined && object.fields !== null) {
            Object.entries(object.fields).forEach(([key, value]) => {
                message.fields[key] = Value.fromJSON(value);
            });
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        obj.fields = {};
        if (message.fields) {
            Object.entries(message.fields).forEach(([k, v]) => {
                obj.fields[k] = Value.toJSON(v);
            });
        }
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseStruct };
        message.fields = {};
        if (object.fields !== undefined && object.fields !== null) {
            Object.entries(object.fields).forEach(([key, value]) => {
                if (value !== undefined) {
                    message.fields[key] = Value.fromPartial(value);
                }
            });
        }
        return message;
    },
};
const baseStruct_FieldsEntry = { key: "" };
export const Struct_FieldsEntry = {
    encode(message, writer = Writer.create()) {
        if (message.key !== "") {
            writer.uint32(10).string(message.key);
        }
        if (message.value !== undefined) {
            Value.encode(message.value, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseStruct_FieldsEntry };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.key = reader.string();
                    break;
                case 2:
                    message.value = Value.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseStruct_FieldsEntry };
        if (object.key !== undefined && object.key !== null) {
            message.key = String(object.key);
        }
        else {
            message.key = "";
        }
        if (object.value !== undefined && object.value !== null) {
            message.value = Value.fromJSON(object.value);
        }
        else {
            message.value = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.key !== undefined && (obj.key = message.key);
        message.value !== undefined &&
            (obj.value = message.value ? Value.toJSON(message.value) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseStruct_FieldsEntry };
        if (object.key !== undefined && object.key !== null) {
            message.key = object.key;
        }
        else {
            message.key = "";
        }
        if (object.value !== undefined && object.value !== null) {
            message.value = Value.fromPartial(object.value);
        }
        else {
            message.value = undefined;
        }
        return message;
    },
};
const baseValue = {};
export const Value = {
    encode(message, writer = Writer.create()) {
        if (message.nullValue !== undefined) {
            writer.uint32(8).int32(message.nullValue);
        }
        if (message.numberValue !== undefined) {
            writer.uint32(17).double(message.numberValue);
        }
        if (message.stringValue !== undefined) {
            writer.uint32(26).string(message.stringValue);
        }
        if (message.boolValue !== undefined) {
            writer.uint32(32).bool(message.boolValue);
        }
        if (message.structValue !== undefined) {
            Struct.encode(message.structValue, writer.uint32(42).fork()).ldelim();
        }
        if (message.listValue !== undefined) {
            ListValue.encode(message.listValue, writer.uint32(50).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseValue };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.nullValue = reader.int32();
                    break;
                case 2:
                    message.numberValue = reader.double();
                    break;
                case 3:
                    message.stringValue = reader.string();
                    break;
                case 4:
                    message.boolValue = reader.bool();
                    break;
                case 5:
                    message.structValue = Struct.decode(reader, reader.uint32());
                    break;
                case 6:
                    message.listValue = ListValue.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseValue };
        if (object.nullValue !== undefined && object.nullValue !== null) {
            message.nullValue = nullValueFromJSON(object.nullValue);
        }
        else {
            message.nullValue = undefined;
        }
        if (object.numberValue !== undefined && object.numberValue !== null) {
            message.numberValue = Number(object.numberValue);
        }
        else {
            message.numberValue = undefined;
        }
        if (object.stringValue !== undefined && object.stringValue !== null) {
            message.stringValue = String(object.stringValue);
        }
        else {
            message.stringValue = undefined;
        }
        if (object.boolValue !== undefined && object.boolValue !== null) {
            message.boolValue = Boolean(object.boolValue);
        }
        else {
            message.boolValue = undefined;
        }
        if (object.structValue !== undefined && object.structValue !== null) {
            message.structValue = Struct.fromJSON(object.structValue);
        }
        else {
            message.structValue = undefined;
        }
        if (object.listValue !== undefined && object.listValue !== null) {
            message.listValue = ListValue.fromJSON(object.listValue);
        }
        else {
            message.listValue = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.nullValue !== undefined &&
            (obj.nullValue =
                message.nullValue !== undefined
                    ? nullValueToJSON(message.nullValue)
                    : undefined);
        message.numberValue !== undefined &&
            (obj.numberValue = message.numberValue);
        message.stringValue !== undefined &&
            (obj.stringValue = message.stringValue);
        message.boolValue !== undefined && (obj.boolValue = message.boolValue);
        message.structValue !== undefined &&
            (obj.structValue = message.structValue
                ? Struct.toJSON(message.structValue)
                : undefined);
        message.listValue !== undefined &&
            (obj.listValue = message.listValue
                ? ListValue.toJSON(message.listValue)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseValue };
        if (object.nullValue !== undefined && object.nullValue !== null) {
            message.nullValue = object.nullValue;
        }
        else {
            message.nullValue = undefined;
        }
        if (object.numberValue !== undefined && object.numberValue !== null) {
            message.numberValue = object.numberValue;
        }
        else {
            message.numberValue = undefined;
        }
        if (object.stringValue !== undefined && object.stringValue !== null) {
            message.stringValue = object.stringValue;
        }
        else {
            message.stringValue = undefined;
        }
        if (object.boolValue !== undefined && object.boolValue !== null) {
            message.boolValue = object.boolValue;
        }
        else {
            message.boolValue = undefined;
        }
        if (object.structValue !== undefined && object.structValue !== null) {
            message.structValue = Struct.fromPartial(object.structValue);
        }
        else {
            message.structValue = undefined;
        }
        if (object.listValue !== undefined && object.listValue !== null) {
            message.listValue = ListValue.fromPartial(object.listValue);
        }
        else {
            message.listValue = undefined;
        }
        return message;
    },
};
const baseListValue = {};
export const ListValue = {
    encode(message, writer = Writer.create()) {
        for (const v of message.values) {
            Value.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseListValue };
        message.values = [];
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.values.push(Value.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseListValue };
        message.values = [];
        if (object.values !== undefined && object.values !== null) {
            for (const e of object.values) {
                message.values.push(Value.fromJSON(e));
            }
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        if (message.values) {
            obj.values = message.values.map((e) => (e ? Value.toJSON(e) : undefined));
        }
        else {
            obj.values = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseListValue };
        message.values = [];
        if (object.values !== undefined && object.values !== null) {
            for (const e of object.values) {
                message.values.push(Value.fromPartial(e));
            }
        }
        return message;
    },
};
