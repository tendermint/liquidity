/* eslint-disable */
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";
import { Value } from "../../google/protobuf/struct";
export const protobufPackage = "grpc.gateway.protoc_gen_openapiv2.options";
/**
 * Scheme describes the schemes supported by the OpenAPI Swagger
 * and Operation objects.
 */
export var Scheme;
(function (Scheme) {
    Scheme[Scheme["UNKNOWN"] = 0] = "UNKNOWN";
    Scheme[Scheme["HTTP"] = 1] = "HTTP";
    Scheme[Scheme["HTTPS"] = 2] = "HTTPS";
    Scheme[Scheme["WS"] = 3] = "WS";
    Scheme[Scheme["WSS"] = 4] = "WSS";
    Scheme[Scheme["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(Scheme || (Scheme = {}));
export function schemeFromJSON(object) {
    switch (object) {
        case 0:
        case "UNKNOWN":
            return Scheme.UNKNOWN;
        case 1:
        case "HTTP":
            return Scheme.HTTP;
        case 2:
        case "HTTPS":
            return Scheme.HTTPS;
        case 3:
        case "WS":
            return Scheme.WS;
        case 4:
        case "WSS":
            return Scheme.WSS;
        case -1:
        case "UNRECOGNIZED":
        default:
            return Scheme.UNRECOGNIZED;
    }
}
export function schemeToJSON(object) {
    switch (object) {
        case Scheme.UNKNOWN:
            return "UNKNOWN";
        case Scheme.HTTP:
            return "HTTP";
        case Scheme.HTTPS:
            return "HTTPS";
        case Scheme.WS:
            return "WS";
        case Scheme.WSS:
            return "WSS";
        default:
            return "UNKNOWN";
    }
}
export var JSONSchema_JSONSchemaSimpleTypes;
(function (JSONSchema_JSONSchemaSimpleTypes) {
    JSONSchema_JSONSchemaSimpleTypes[JSONSchema_JSONSchemaSimpleTypes["UNKNOWN"] = 0] = "UNKNOWN";
    JSONSchema_JSONSchemaSimpleTypes[JSONSchema_JSONSchemaSimpleTypes["ARRAY"] = 1] = "ARRAY";
    JSONSchema_JSONSchemaSimpleTypes[JSONSchema_JSONSchemaSimpleTypes["BOOLEAN"] = 2] = "BOOLEAN";
    JSONSchema_JSONSchemaSimpleTypes[JSONSchema_JSONSchemaSimpleTypes["INTEGER"] = 3] = "INTEGER";
    JSONSchema_JSONSchemaSimpleTypes[JSONSchema_JSONSchemaSimpleTypes["NULL"] = 4] = "NULL";
    JSONSchema_JSONSchemaSimpleTypes[JSONSchema_JSONSchemaSimpleTypes["NUMBER"] = 5] = "NUMBER";
    JSONSchema_JSONSchemaSimpleTypes[JSONSchema_JSONSchemaSimpleTypes["OBJECT"] = 6] = "OBJECT";
    JSONSchema_JSONSchemaSimpleTypes[JSONSchema_JSONSchemaSimpleTypes["STRING"] = 7] = "STRING";
    JSONSchema_JSONSchemaSimpleTypes[JSONSchema_JSONSchemaSimpleTypes["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(JSONSchema_JSONSchemaSimpleTypes || (JSONSchema_JSONSchemaSimpleTypes = {}));
export function jSONSchema_JSONSchemaSimpleTypesFromJSON(object) {
    switch (object) {
        case 0:
        case "UNKNOWN":
            return JSONSchema_JSONSchemaSimpleTypes.UNKNOWN;
        case 1:
        case "ARRAY":
            return JSONSchema_JSONSchemaSimpleTypes.ARRAY;
        case 2:
        case "BOOLEAN":
            return JSONSchema_JSONSchemaSimpleTypes.BOOLEAN;
        case 3:
        case "INTEGER":
            return JSONSchema_JSONSchemaSimpleTypes.INTEGER;
        case 4:
        case "NULL":
            return JSONSchema_JSONSchemaSimpleTypes.NULL;
        case 5:
        case "NUMBER":
            return JSONSchema_JSONSchemaSimpleTypes.NUMBER;
        case 6:
        case "OBJECT":
            return JSONSchema_JSONSchemaSimpleTypes.OBJECT;
        case 7:
        case "STRING":
            return JSONSchema_JSONSchemaSimpleTypes.STRING;
        case -1:
        case "UNRECOGNIZED":
        default:
            return JSONSchema_JSONSchemaSimpleTypes.UNRECOGNIZED;
    }
}
export function jSONSchema_JSONSchemaSimpleTypesToJSON(object) {
    switch (object) {
        case JSONSchema_JSONSchemaSimpleTypes.UNKNOWN:
            return "UNKNOWN";
        case JSONSchema_JSONSchemaSimpleTypes.ARRAY:
            return "ARRAY";
        case JSONSchema_JSONSchemaSimpleTypes.BOOLEAN:
            return "BOOLEAN";
        case JSONSchema_JSONSchemaSimpleTypes.INTEGER:
            return "INTEGER";
        case JSONSchema_JSONSchemaSimpleTypes.NULL:
            return "NULL";
        case JSONSchema_JSONSchemaSimpleTypes.NUMBER:
            return "NUMBER";
        case JSONSchema_JSONSchemaSimpleTypes.OBJECT:
            return "OBJECT";
        case JSONSchema_JSONSchemaSimpleTypes.STRING:
            return "STRING";
        default:
            return "UNKNOWN";
    }
}
/**
 * The type of the security scheme. Valid values are "basic",
 * "apiKey" or "oauth2".
 */
export var SecurityScheme_Type;
(function (SecurityScheme_Type) {
    SecurityScheme_Type[SecurityScheme_Type["TYPE_INVALID"] = 0] = "TYPE_INVALID";
    SecurityScheme_Type[SecurityScheme_Type["TYPE_BASIC"] = 1] = "TYPE_BASIC";
    SecurityScheme_Type[SecurityScheme_Type["TYPE_API_KEY"] = 2] = "TYPE_API_KEY";
    SecurityScheme_Type[SecurityScheme_Type["TYPE_OAUTH2"] = 3] = "TYPE_OAUTH2";
    SecurityScheme_Type[SecurityScheme_Type["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(SecurityScheme_Type || (SecurityScheme_Type = {}));
export function securityScheme_TypeFromJSON(object) {
    switch (object) {
        case 0:
        case "TYPE_INVALID":
            return SecurityScheme_Type.TYPE_INVALID;
        case 1:
        case "TYPE_BASIC":
            return SecurityScheme_Type.TYPE_BASIC;
        case 2:
        case "TYPE_API_KEY":
            return SecurityScheme_Type.TYPE_API_KEY;
        case 3:
        case "TYPE_OAUTH2":
            return SecurityScheme_Type.TYPE_OAUTH2;
        case -1:
        case "UNRECOGNIZED":
        default:
            return SecurityScheme_Type.UNRECOGNIZED;
    }
}
export function securityScheme_TypeToJSON(object) {
    switch (object) {
        case SecurityScheme_Type.TYPE_INVALID:
            return "TYPE_INVALID";
        case SecurityScheme_Type.TYPE_BASIC:
            return "TYPE_BASIC";
        case SecurityScheme_Type.TYPE_API_KEY:
            return "TYPE_API_KEY";
        case SecurityScheme_Type.TYPE_OAUTH2:
            return "TYPE_OAUTH2";
        default:
            return "UNKNOWN";
    }
}
/** The location of the API key. Valid values are "query" or "header". */
export var SecurityScheme_In;
(function (SecurityScheme_In) {
    SecurityScheme_In[SecurityScheme_In["IN_INVALID"] = 0] = "IN_INVALID";
    SecurityScheme_In[SecurityScheme_In["IN_QUERY"] = 1] = "IN_QUERY";
    SecurityScheme_In[SecurityScheme_In["IN_HEADER"] = 2] = "IN_HEADER";
    SecurityScheme_In[SecurityScheme_In["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(SecurityScheme_In || (SecurityScheme_In = {}));
export function securityScheme_InFromJSON(object) {
    switch (object) {
        case 0:
        case "IN_INVALID":
            return SecurityScheme_In.IN_INVALID;
        case 1:
        case "IN_QUERY":
            return SecurityScheme_In.IN_QUERY;
        case 2:
        case "IN_HEADER":
            return SecurityScheme_In.IN_HEADER;
        case -1:
        case "UNRECOGNIZED":
        default:
            return SecurityScheme_In.UNRECOGNIZED;
    }
}
export function securityScheme_InToJSON(object) {
    switch (object) {
        case SecurityScheme_In.IN_INVALID:
            return "IN_INVALID";
        case SecurityScheme_In.IN_QUERY:
            return "IN_QUERY";
        case SecurityScheme_In.IN_HEADER:
            return "IN_HEADER";
        default:
            return "UNKNOWN";
    }
}
/**
 * The flow used by the OAuth2 security scheme. Valid values are
 * "implicit", "password", "application" or "accessCode".
 */
export var SecurityScheme_Flow;
(function (SecurityScheme_Flow) {
    SecurityScheme_Flow[SecurityScheme_Flow["FLOW_INVALID"] = 0] = "FLOW_INVALID";
    SecurityScheme_Flow[SecurityScheme_Flow["FLOW_IMPLICIT"] = 1] = "FLOW_IMPLICIT";
    SecurityScheme_Flow[SecurityScheme_Flow["FLOW_PASSWORD"] = 2] = "FLOW_PASSWORD";
    SecurityScheme_Flow[SecurityScheme_Flow["FLOW_APPLICATION"] = 3] = "FLOW_APPLICATION";
    SecurityScheme_Flow[SecurityScheme_Flow["FLOW_ACCESS_CODE"] = 4] = "FLOW_ACCESS_CODE";
    SecurityScheme_Flow[SecurityScheme_Flow["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(SecurityScheme_Flow || (SecurityScheme_Flow = {}));
export function securityScheme_FlowFromJSON(object) {
    switch (object) {
        case 0:
        case "FLOW_INVALID":
            return SecurityScheme_Flow.FLOW_INVALID;
        case 1:
        case "FLOW_IMPLICIT":
            return SecurityScheme_Flow.FLOW_IMPLICIT;
        case 2:
        case "FLOW_PASSWORD":
            return SecurityScheme_Flow.FLOW_PASSWORD;
        case 3:
        case "FLOW_APPLICATION":
            return SecurityScheme_Flow.FLOW_APPLICATION;
        case 4:
        case "FLOW_ACCESS_CODE":
            return SecurityScheme_Flow.FLOW_ACCESS_CODE;
        case -1:
        case "UNRECOGNIZED":
        default:
            return SecurityScheme_Flow.UNRECOGNIZED;
    }
}
export function securityScheme_FlowToJSON(object) {
    switch (object) {
        case SecurityScheme_Flow.FLOW_INVALID:
            return "FLOW_INVALID";
        case SecurityScheme_Flow.FLOW_IMPLICIT:
            return "FLOW_IMPLICIT";
        case SecurityScheme_Flow.FLOW_PASSWORD:
            return "FLOW_PASSWORD";
        case SecurityScheme_Flow.FLOW_APPLICATION:
            return "FLOW_APPLICATION";
        case SecurityScheme_Flow.FLOW_ACCESS_CODE:
            return "FLOW_ACCESS_CODE";
        default:
            return "UNKNOWN";
    }
}
const baseSwagger = {
    swagger: "",
    host: "",
    basePath: "",
    schemes: 0,
    consumes: "",
    produces: "",
};
export const Swagger = {
    encode(message, writer = Writer.create()) {
        if (message.swagger !== "") {
            writer.uint32(10).string(message.swagger);
        }
        if (message.info !== undefined) {
            Info.encode(message.info, writer.uint32(18).fork()).ldelim();
        }
        if (message.host !== "") {
            writer.uint32(26).string(message.host);
        }
        if (message.basePath !== "") {
            writer.uint32(34).string(message.basePath);
        }
        writer.uint32(42).fork();
        for (const v of message.schemes) {
            writer.int32(v);
        }
        writer.ldelim();
        for (const v of message.consumes) {
            writer.uint32(50).string(v);
        }
        for (const v of message.produces) {
            writer.uint32(58).string(v);
        }
        Object.entries(message.responses).forEach(([key, value]) => {
            Swagger_ResponsesEntry.encode({ key: key, value }, writer.uint32(82).fork()).ldelim();
        });
        if (message.securityDefinitions !== undefined) {
            SecurityDefinitions.encode(message.securityDefinitions, writer.uint32(90).fork()).ldelim();
        }
        for (const v of message.security) {
            SecurityRequirement.encode(v, writer.uint32(98).fork()).ldelim();
        }
        if (message.externalDocs !== undefined) {
            ExternalDocumentation.encode(message.externalDocs, writer.uint32(114).fork()).ldelim();
        }
        Object.entries(message.extensions).forEach(([key, value]) => {
            Swagger_ExtensionsEntry.encode({ key: key, value }, writer.uint32(122).fork()).ldelim();
        });
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseSwagger };
        message.schemes = [];
        message.consumes = [];
        message.produces = [];
        message.responses = {};
        message.security = [];
        message.extensions = {};
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.swagger = reader.string();
                    break;
                case 2:
                    message.info = Info.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.host = reader.string();
                    break;
                case 4:
                    message.basePath = reader.string();
                    break;
                case 5:
                    if ((tag & 7) === 2) {
                        const end2 = reader.uint32() + reader.pos;
                        while (reader.pos < end2) {
                            message.schemes.push(reader.int32());
                        }
                    }
                    else {
                        message.schemes.push(reader.int32());
                    }
                    break;
                case 6:
                    message.consumes.push(reader.string());
                    break;
                case 7:
                    message.produces.push(reader.string());
                    break;
                case 10:
                    const entry10 = Swagger_ResponsesEntry.decode(reader, reader.uint32());
                    if (entry10.value !== undefined) {
                        message.responses[entry10.key] = entry10.value;
                    }
                    break;
                case 11:
                    message.securityDefinitions = SecurityDefinitions.decode(reader, reader.uint32());
                    break;
                case 12:
                    message.security.push(SecurityRequirement.decode(reader, reader.uint32()));
                    break;
                case 14:
                    message.externalDocs = ExternalDocumentation.decode(reader, reader.uint32());
                    break;
                case 15:
                    const entry15 = Swagger_ExtensionsEntry.decode(reader, reader.uint32());
                    if (entry15.value !== undefined) {
                        message.extensions[entry15.key] = entry15.value;
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
        const message = { ...baseSwagger };
        message.schemes = [];
        message.consumes = [];
        message.produces = [];
        message.responses = {};
        message.security = [];
        message.extensions = {};
        if (object.swagger !== undefined && object.swagger !== null) {
            message.swagger = String(object.swagger);
        }
        else {
            message.swagger = "";
        }
        if (object.info !== undefined && object.info !== null) {
            message.info = Info.fromJSON(object.info);
        }
        else {
            message.info = undefined;
        }
        if (object.host !== undefined && object.host !== null) {
            message.host = String(object.host);
        }
        else {
            message.host = "";
        }
        if (object.basePath !== undefined && object.basePath !== null) {
            message.basePath = String(object.basePath);
        }
        else {
            message.basePath = "";
        }
        if (object.schemes !== undefined && object.schemes !== null) {
            for (const e of object.schemes) {
                message.schemes.push(schemeFromJSON(e));
            }
        }
        if (object.consumes !== undefined && object.consumes !== null) {
            for (const e of object.consumes) {
                message.consumes.push(String(e));
            }
        }
        if (object.produces !== undefined && object.produces !== null) {
            for (const e of object.produces) {
                message.produces.push(String(e));
            }
        }
        if (object.responses !== undefined && object.responses !== null) {
            Object.entries(object.responses).forEach(([key, value]) => {
                message.responses[key] = Response.fromJSON(value);
            });
        }
        if (object.securityDefinitions !== undefined &&
            object.securityDefinitions !== null) {
            message.securityDefinitions = SecurityDefinitions.fromJSON(object.securityDefinitions);
        }
        else {
            message.securityDefinitions = undefined;
        }
        if (object.security !== undefined && object.security !== null) {
            for (const e of object.security) {
                message.security.push(SecurityRequirement.fromJSON(e));
            }
        }
        if (object.externalDocs !== undefined && object.externalDocs !== null) {
            message.externalDocs = ExternalDocumentation.fromJSON(object.externalDocs);
        }
        else {
            message.externalDocs = undefined;
        }
        if (object.extensions !== undefined && object.extensions !== null) {
            Object.entries(object.extensions).forEach(([key, value]) => {
                message.extensions[key] = Value.fromJSON(value);
            });
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.swagger !== undefined && (obj.swagger = message.swagger);
        message.info !== undefined &&
            (obj.info = message.info ? Info.toJSON(message.info) : undefined);
        message.host !== undefined && (obj.host = message.host);
        message.basePath !== undefined && (obj.basePath = message.basePath);
        if (message.schemes) {
            obj.schemes = message.schemes.map((e) => schemeToJSON(e));
        }
        else {
            obj.schemes = [];
        }
        if (message.consumes) {
            obj.consumes = message.consumes.map((e) => e);
        }
        else {
            obj.consumes = [];
        }
        if (message.produces) {
            obj.produces = message.produces.map((e) => e);
        }
        else {
            obj.produces = [];
        }
        obj.responses = {};
        if (message.responses) {
            Object.entries(message.responses).forEach(([k, v]) => {
                obj.responses[k] = Response.toJSON(v);
            });
        }
        message.securityDefinitions !== undefined &&
            (obj.securityDefinitions = message.securityDefinitions
                ? SecurityDefinitions.toJSON(message.securityDefinitions)
                : undefined);
        if (message.security) {
            obj.security = message.security.map((e) => e ? SecurityRequirement.toJSON(e) : undefined);
        }
        else {
            obj.security = [];
        }
        message.externalDocs !== undefined &&
            (obj.externalDocs = message.externalDocs
                ? ExternalDocumentation.toJSON(message.externalDocs)
                : undefined);
        obj.extensions = {};
        if (message.extensions) {
            Object.entries(message.extensions).forEach(([k, v]) => {
                obj.extensions[k] = Value.toJSON(v);
            });
        }
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseSwagger };
        message.schemes = [];
        message.consumes = [];
        message.produces = [];
        message.responses = {};
        message.security = [];
        message.extensions = {};
        if (object.swagger !== undefined && object.swagger !== null) {
            message.swagger = object.swagger;
        }
        else {
            message.swagger = "";
        }
        if (object.info !== undefined && object.info !== null) {
            message.info = Info.fromPartial(object.info);
        }
        else {
            message.info = undefined;
        }
        if (object.host !== undefined && object.host !== null) {
            message.host = object.host;
        }
        else {
            message.host = "";
        }
        if (object.basePath !== undefined && object.basePath !== null) {
            message.basePath = object.basePath;
        }
        else {
            message.basePath = "";
        }
        if (object.schemes !== undefined && object.schemes !== null) {
            for (const e of object.schemes) {
                message.schemes.push(e);
            }
        }
        if (object.consumes !== undefined && object.consumes !== null) {
            for (const e of object.consumes) {
                message.consumes.push(e);
            }
        }
        if (object.produces !== undefined && object.produces !== null) {
            for (const e of object.produces) {
                message.produces.push(e);
            }
        }
        if (object.responses !== undefined && object.responses !== null) {
            Object.entries(object.responses).forEach(([key, value]) => {
                if (value !== undefined) {
                    message.responses[key] = Response.fromPartial(value);
                }
            });
        }
        if (object.securityDefinitions !== undefined &&
            object.securityDefinitions !== null) {
            message.securityDefinitions = SecurityDefinitions.fromPartial(object.securityDefinitions);
        }
        else {
            message.securityDefinitions = undefined;
        }
        if (object.security !== undefined && object.security !== null) {
            for (const e of object.security) {
                message.security.push(SecurityRequirement.fromPartial(e));
            }
        }
        if (object.externalDocs !== undefined && object.externalDocs !== null) {
            message.externalDocs = ExternalDocumentation.fromPartial(object.externalDocs);
        }
        else {
            message.externalDocs = undefined;
        }
        if (object.extensions !== undefined && object.extensions !== null) {
            Object.entries(object.extensions).forEach(([key, value]) => {
                if (value !== undefined) {
                    message.extensions[key] = Value.fromPartial(value);
                }
            });
        }
        return message;
    },
};
const baseSwagger_ResponsesEntry = { key: "" };
export const Swagger_ResponsesEntry = {
    encode(message, writer = Writer.create()) {
        if (message.key !== "") {
            writer.uint32(10).string(message.key);
        }
        if (message.value !== undefined) {
            Response.encode(message.value, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseSwagger_ResponsesEntry };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.key = reader.string();
                    break;
                case 2:
                    message.value = Response.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseSwagger_ResponsesEntry };
        if (object.key !== undefined && object.key !== null) {
            message.key = String(object.key);
        }
        else {
            message.key = "";
        }
        if (object.value !== undefined && object.value !== null) {
            message.value = Response.fromJSON(object.value);
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
            (obj.value = message.value ? Response.toJSON(message.value) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseSwagger_ResponsesEntry };
        if (object.key !== undefined && object.key !== null) {
            message.key = object.key;
        }
        else {
            message.key = "";
        }
        if (object.value !== undefined && object.value !== null) {
            message.value = Response.fromPartial(object.value);
        }
        else {
            message.value = undefined;
        }
        return message;
    },
};
const baseSwagger_ExtensionsEntry = { key: "" };
export const Swagger_ExtensionsEntry = {
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
        const message = {
            ...baseSwagger_ExtensionsEntry,
        };
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
        const message = {
            ...baseSwagger_ExtensionsEntry,
        };
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
        const message = {
            ...baseSwagger_ExtensionsEntry,
        };
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
const baseOperation = {
    tags: "",
    summary: "",
    description: "",
    operationId: "",
    consumes: "",
    produces: "",
    schemes: 0,
    deprecated: false,
};
export const Operation = {
    encode(message, writer = Writer.create()) {
        for (const v of message.tags) {
            writer.uint32(10).string(v);
        }
        if (message.summary !== "") {
            writer.uint32(18).string(message.summary);
        }
        if (message.description !== "") {
            writer.uint32(26).string(message.description);
        }
        if (message.externalDocs !== undefined) {
            ExternalDocumentation.encode(message.externalDocs, writer.uint32(34).fork()).ldelim();
        }
        if (message.operationId !== "") {
            writer.uint32(42).string(message.operationId);
        }
        for (const v of message.consumes) {
            writer.uint32(50).string(v);
        }
        for (const v of message.produces) {
            writer.uint32(58).string(v);
        }
        Object.entries(message.responses).forEach(([key, value]) => {
            Operation_ResponsesEntry.encode({ key: key, value }, writer.uint32(74).fork()).ldelim();
        });
        writer.uint32(82).fork();
        for (const v of message.schemes) {
            writer.int32(v);
        }
        writer.ldelim();
        if (message.deprecated === true) {
            writer.uint32(88).bool(message.deprecated);
        }
        for (const v of message.security) {
            SecurityRequirement.encode(v, writer.uint32(98).fork()).ldelim();
        }
        Object.entries(message.extensions).forEach(([key, value]) => {
            Operation_ExtensionsEntry.encode({ key: key, value }, writer.uint32(106).fork()).ldelim();
        });
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseOperation };
        message.tags = [];
        message.consumes = [];
        message.produces = [];
        message.responses = {};
        message.schemes = [];
        message.security = [];
        message.extensions = {};
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.tags.push(reader.string());
                    break;
                case 2:
                    message.summary = reader.string();
                    break;
                case 3:
                    message.description = reader.string();
                    break;
                case 4:
                    message.externalDocs = ExternalDocumentation.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.operationId = reader.string();
                    break;
                case 6:
                    message.consumes.push(reader.string());
                    break;
                case 7:
                    message.produces.push(reader.string());
                    break;
                case 9:
                    const entry9 = Operation_ResponsesEntry.decode(reader, reader.uint32());
                    if (entry9.value !== undefined) {
                        message.responses[entry9.key] = entry9.value;
                    }
                    break;
                case 10:
                    if ((tag & 7) === 2) {
                        const end2 = reader.uint32() + reader.pos;
                        while (reader.pos < end2) {
                            message.schemes.push(reader.int32());
                        }
                    }
                    else {
                        message.schemes.push(reader.int32());
                    }
                    break;
                case 11:
                    message.deprecated = reader.bool();
                    break;
                case 12:
                    message.security.push(SecurityRequirement.decode(reader, reader.uint32()));
                    break;
                case 13:
                    const entry13 = Operation_ExtensionsEntry.decode(reader, reader.uint32());
                    if (entry13.value !== undefined) {
                        message.extensions[entry13.key] = entry13.value;
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
        const message = { ...baseOperation };
        message.tags = [];
        message.consumes = [];
        message.produces = [];
        message.responses = {};
        message.schemes = [];
        message.security = [];
        message.extensions = {};
        if (object.tags !== undefined && object.tags !== null) {
            for (const e of object.tags) {
                message.tags.push(String(e));
            }
        }
        if (object.summary !== undefined && object.summary !== null) {
            message.summary = String(object.summary);
        }
        else {
            message.summary = "";
        }
        if (object.description !== undefined && object.description !== null) {
            message.description = String(object.description);
        }
        else {
            message.description = "";
        }
        if (object.externalDocs !== undefined && object.externalDocs !== null) {
            message.externalDocs = ExternalDocumentation.fromJSON(object.externalDocs);
        }
        else {
            message.externalDocs = undefined;
        }
        if (object.operationId !== undefined && object.operationId !== null) {
            message.operationId = String(object.operationId);
        }
        else {
            message.operationId = "";
        }
        if (object.consumes !== undefined && object.consumes !== null) {
            for (const e of object.consumes) {
                message.consumes.push(String(e));
            }
        }
        if (object.produces !== undefined && object.produces !== null) {
            for (const e of object.produces) {
                message.produces.push(String(e));
            }
        }
        if (object.responses !== undefined && object.responses !== null) {
            Object.entries(object.responses).forEach(([key, value]) => {
                message.responses[key] = Response.fromJSON(value);
            });
        }
        if (object.schemes !== undefined && object.schemes !== null) {
            for (const e of object.schemes) {
                message.schemes.push(schemeFromJSON(e));
            }
        }
        if (object.deprecated !== undefined && object.deprecated !== null) {
            message.deprecated = Boolean(object.deprecated);
        }
        else {
            message.deprecated = false;
        }
        if (object.security !== undefined && object.security !== null) {
            for (const e of object.security) {
                message.security.push(SecurityRequirement.fromJSON(e));
            }
        }
        if (object.extensions !== undefined && object.extensions !== null) {
            Object.entries(object.extensions).forEach(([key, value]) => {
                message.extensions[key] = Value.fromJSON(value);
            });
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        if (message.tags) {
            obj.tags = message.tags.map((e) => e);
        }
        else {
            obj.tags = [];
        }
        message.summary !== undefined && (obj.summary = message.summary);
        message.description !== undefined &&
            (obj.description = message.description);
        message.externalDocs !== undefined &&
            (obj.externalDocs = message.externalDocs
                ? ExternalDocumentation.toJSON(message.externalDocs)
                : undefined);
        message.operationId !== undefined &&
            (obj.operationId = message.operationId);
        if (message.consumes) {
            obj.consumes = message.consumes.map((e) => e);
        }
        else {
            obj.consumes = [];
        }
        if (message.produces) {
            obj.produces = message.produces.map((e) => e);
        }
        else {
            obj.produces = [];
        }
        obj.responses = {};
        if (message.responses) {
            Object.entries(message.responses).forEach(([k, v]) => {
                obj.responses[k] = Response.toJSON(v);
            });
        }
        if (message.schemes) {
            obj.schemes = message.schemes.map((e) => schemeToJSON(e));
        }
        else {
            obj.schemes = [];
        }
        message.deprecated !== undefined && (obj.deprecated = message.deprecated);
        if (message.security) {
            obj.security = message.security.map((e) => e ? SecurityRequirement.toJSON(e) : undefined);
        }
        else {
            obj.security = [];
        }
        obj.extensions = {};
        if (message.extensions) {
            Object.entries(message.extensions).forEach(([k, v]) => {
                obj.extensions[k] = Value.toJSON(v);
            });
        }
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseOperation };
        message.tags = [];
        message.consumes = [];
        message.produces = [];
        message.responses = {};
        message.schemes = [];
        message.security = [];
        message.extensions = {};
        if (object.tags !== undefined && object.tags !== null) {
            for (const e of object.tags) {
                message.tags.push(e);
            }
        }
        if (object.summary !== undefined && object.summary !== null) {
            message.summary = object.summary;
        }
        else {
            message.summary = "";
        }
        if (object.description !== undefined && object.description !== null) {
            message.description = object.description;
        }
        else {
            message.description = "";
        }
        if (object.externalDocs !== undefined && object.externalDocs !== null) {
            message.externalDocs = ExternalDocumentation.fromPartial(object.externalDocs);
        }
        else {
            message.externalDocs = undefined;
        }
        if (object.operationId !== undefined && object.operationId !== null) {
            message.operationId = object.operationId;
        }
        else {
            message.operationId = "";
        }
        if (object.consumes !== undefined && object.consumes !== null) {
            for (const e of object.consumes) {
                message.consumes.push(e);
            }
        }
        if (object.produces !== undefined && object.produces !== null) {
            for (const e of object.produces) {
                message.produces.push(e);
            }
        }
        if (object.responses !== undefined && object.responses !== null) {
            Object.entries(object.responses).forEach(([key, value]) => {
                if (value !== undefined) {
                    message.responses[key] = Response.fromPartial(value);
                }
            });
        }
        if (object.schemes !== undefined && object.schemes !== null) {
            for (const e of object.schemes) {
                message.schemes.push(e);
            }
        }
        if (object.deprecated !== undefined && object.deprecated !== null) {
            message.deprecated = object.deprecated;
        }
        else {
            message.deprecated = false;
        }
        if (object.security !== undefined && object.security !== null) {
            for (const e of object.security) {
                message.security.push(SecurityRequirement.fromPartial(e));
            }
        }
        if (object.extensions !== undefined && object.extensions !== null) {
            Object.entries(object.extensions).forEach(([key, value]) => {
                if (value !== undefined) {
                    message.extensions[key] = Value.fromPartial(value);
                }
            });
        }
        return message;
    },
};
const baseOperation_ResponsesEntry = { key: "" };
export const Operation_ResponsesEntry = {
    encode(message, writer = Writer.create()) {
        if (message.key !== "") {
            writer.uint32(10).string(message.key);
        }
        if (message.value !== undefined) {
            Response.encode(message.value, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseOperation_ResponsesEntry,
        };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.key = reader.string();
                    break;
                case 2:
                    message.value = Response.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseOperation_ResponsesEntry,
        };
        if (object.key !== undefined && object.key !== null) {
            message.key = String(object.key);
        }
        else {
            message.key = "";
        }
        if (object.value !== undefined && object.value !== null) {
            message.value = Response.fromJSON(object.value);
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
            (obj.value = message.value ? Response.toJSON(message.value) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseOperation_ResponsesEntry,
        };
        if (object.key !== undefined && object.key !== null) {
            message.key = object.key;
        }
        else {
            message.key = "";
        }
        if (object.value !== undefined && object.value !== null) {
            message.value = Response.fromPartial(object.value);
        }
        else {
            message.value = undefined;
        }
        return message;
    },
};
const baseOperation_ExtensionsEntry = { key: "" };
export const Operation_ExtensionsEntry = {
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
        const message = {
            ...baseOperation_ExtensionsEntry,
        };
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
        const message = {
            ...baseOperation_ExtensionsEntry,
        };
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
        const message = {
            ...baseOperation_ExtensionsEntry,
        };
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
const baseHeader = {
    description: "",
    type: "",
    format: "",
    default: "",
    pattern: "",
};
export const Header = {
    encode(message, writer = Writer.create()) {
        if (message.description !== "") {
            writer.uint32(10).string(message.description);
        }
        if (message.type !== "") {
            writer.uint32(18).string(message.type);
        }
        if (message.format !== "") {
            writer.uint32(26).string(message.format);
        }
        if (message.default !== "") {
            writer.uint32(50).string(message.default);
        }
        if (message.pattern !== "") {
            writer.uint32(106).string(message.pattern);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseHeader };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.description = reader.string();
                    break;
                case 2:
                    message.type = reader.string();
                    break;
                case 3:
                    message.format = reader.string();
                    break;
                case 6:
                    message.default = reader.string();
                    break;
                case 13:
                    message.pattern = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseHeader };
        if (object.description !== undefined && object.description !== null) {
            message.description = String(object.description);
        }
        else {
            message.description = "";
        }
        if (object.type !== undefined && object.type !== null) {
            message.type = String(object.type);
        }
        else {
            message.type = "";
        }
        if (object.format !== undefined && object.format !== null) {
            message.format = String(object.format);
        }
        else {
            message.format = "";
        }
        if (object.default !== undefined && object.default !== null) {
            message.default = String(object.default);
        }
        else {
            message.default = "";
        }
        if (object.pattern !== undefined && object.pattern !== null) {
            message.pattern = String(object.pattern);
        }
        else {
            message.pattern = "";
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.description !== undefined &&
            (obj.description = message.description);
        message.type !== undefined && (obj.type = message.type);
        message.format !== undefined && (obj.format = message.format);
        message.default !== undefined && (obj.default = message.default);
        message.pattern !== undefined && (obj.pattern = message.pattern);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseHeader };
        if (object.description !== undefined && object.description !== null) {
            message.description = object.description;
        }
        else {
            message.description = "";
        }
        if (object.type !== undefined && object.type !== null) {
            message.type = object.type;
        }
        else {
            message.type = "";
        }
        if (object.format !== undefined && object.format !== null) {
            message.format = object.format;
        }
        else {
            message.format = "";
        }
        if (object.default !== undefined && object.default !== null) {
            message.default = object.default;
        }
        else {
            message.default = "";
        }
        if (object.pattern !== undefined && object.pattern !== null) {
            message.pattern = object.pattern;
        }
        else {
            message.pattern = "";
        }
        return message;
    },
};
const baseResponse = { description: "" };
export const Response = {
    encode(message, writer = Writer.create()) {
        if (message.description !== "") {
            writer.uint32(10).string(message.description);
        }
        if (message.schema !== undefined) {
            Schema.encode(message.schema, writer.uint32(18).fork()).ldelim();
        }
        Object.entries(message.headers).forEach(([key, value]) => {
            Response_HeadersEntry.encode({ key: key, value }, writer.uint32(26).fork()).ldelim();
        });
        Object.entries(message.examples).forEach(([key, value]) => {
            Response_ExamplesEntry.encode({ key: key, value }, writer.uint32(34).fork()).ldelim();
        });
        Object.entries(message.extensions).forEach(([key, value]) => {
            Response_ExtensionsEntry.encode({ key: key, value }, writer.uint32(42).fork()).ldelim();
        });
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseResponse };
        message.headers = {};
        message.examples = {};
        message.extensions = {};
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.description = reader.string();
                    break;
                case 2:
                    message.schema = Schema.decode(reader, reader.uint32());
                    break;
                case 3:
                    const entry3 = Response_HeadersEntry.decode(reader, reader.uint32());
                    if (entry3.value !== undefined) {
                        message.headers[entry3.key] = entry3.value;
                    }
                    break;
                case 4:
                    const entry4 = Response_ExamplesEntry.decode(reader, reader.uint32());
                    if (entry4.value !== undefined) {
                        message.examples[entry4.key] = entry4.value;
                    }
                    break;
                case 5:
                    const entry5 = Response_ExtensionsEntry.decode(reader, reader.uint32());
                    if (entry5.value !== undefined) {
                        message.extensions[entry5.key] = entry5.value;
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
        const message = { ...baseResponse };
        message.headers = {};
        message.examples = {};
        message.extensions = {};
        if (object.description !== undefined && object.description !== null) {
            message.description = String(object.description);
        }
        else {
            message.description = "";
        }
        if (object.schema !== undefined && object.schema !== null) {
            message.schema = Schema.fromJSON(object.schema);
        }
        else {
            message.schema = undefined;
        }
        if (object.headers !== undefined && object.headers !== null) {
            Object.entries(object.headers).forEach(([key, value]) => {
                message.headers[key] = Header.fromJSON(value);
            });
        }
        if (object.examples !== undefined && object.examples !== null) {
            Object.entries(object.examples).forEach(([key, value]) => {
                message.examples[key] = String(value);
            });
        }
        if (object.extensions !== undefined && object.extensions !== null) {
            Object.entries(object.extensions).forEach(([key, value]) => {
                message.extensions[key] = Value.fromJSON(value);
            });
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.description !== undefined &&
            (obj.description = message.description);
        message.schema !== undefined &&
            (obj.schema = message.schema ? Schema.toJSON(message.schema) : undefined);
        obj.headers = {};
        if (message.headers) {
            Object.entries(message.headers).forEach(([k, v]) => {
                obj.headers[k] = Header.toJSON(v);
            });
        }
        obj.examples = {};
        if (message.examples) {
            Object.entries(message.examples).forEach(([k, v]) => {
                obj.examples[k] = v;
            });
        }
        obj.extensions = {};
        if (message.extensions) {
            Object.entries(message.extensions).forEach(([k, v]) => {
                obj.extensions[k] = Value.toJSON(v);
            });
        }
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseResponse };
        message.headers = {};
        message.examples = {};
        message.extensions = {};
        if (object.description !== undefined && object.description !== null) {
            message.description = object.description;
        }
        else {
            message.description = "";
        }
        if (object.schema !== undefined && object.schema !== null) {
            message.schema = Schema.fromPartial(object.schema);
        }
        else {
            message.schema = undefined;
        }
        if (object.headers !== undefined && object.headers !== null) {
            Object.entries(object.headers).forEach(([key, value]) => {
                if (value !== undefined) {
                    message.headers[key] = Header.fromPartial(value);
                }
            });
        }
        if (object.examples !== undefined && object.examples !== null) {
            Object.entries(object.examples).forEach(([key, value]) => {
                if (value !== undefined) {
                    message.examples[key] = String(value);
                }
            });
        }
        if (object.extensions !== undefined && object.extensions !== null) {
            Object.entries(object.extensions).forEach(([key, value]) => {
                if (value !== undefined) {
                    message.extensions[key] = Value.fromPartial(value);
                }
            });
        }
        return message;
    },
};
const baseResponse_HeadersEntry = { key: "" };
export const Response_HeadersEntry = {
    encode(message, writer = Writer.create()) {
        if (message.key !== "") {
            writer.uint32(10).string(message.key);
        }
        if (message.value !== undefined) {
            Header.encode(message.value, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseResponse_HeadersEntry };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.key = reader.string();
                    break;
                case 2:
                    message.value = Header.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseResponse_HeadersEntry };
        if (object.key !== undefined && object.key !== null) {
            message.key = String(object.key);
        }
        else {
            message.key = "";
        }
        if (object.value !== undefined && object.value !== null) {
            message.value = Header.fromJSON(object.value);
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
            (obj.value = message.value ? Header.toJSON(message.value) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseResponse_HeadersEntry };
        if (object.key !== undefined && object.key !== null) {
            message.key = object.key;
        }
        else {
            message.key = "";
        }
        if (object.value !== undefined && object.value !== null) {
            message.value = Header.fromPartial(object.value);
        }
        else {
            message.value = undefined;
        }
        return message;
    },
};
const baseResponse_ExamplesEntry = { key: "", value: "" };
export const Response_ExamplesEntry = {
    encode(message, writer = Writer.create()) {
        if (message.key !== "") {
            writer.uint32(10).string(message.key);
        }
        if (message.value !== "") {
            writer.uint32(18).string(message.value);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseResponse_ExamplesEntry };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.key = reader.string();
                    break;
                case 2:
                    message.value = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseResponse_ExamplesEntry };
        if (object.key !== undefined && object.key !== null) {
            message.key = String(object.key);
        }
        else {
            message.key = "";
        }
        if (object.value !== undefined && object.value !== null) {
            message.value = String(object.value);
        }
        else {
            message.value = "";
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.key !== undefined && (obj.key = message.key);
        message.value !== undefined && (obj.value = message.value);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseResponse_ExamplesEntry };
        if (object.key !== undefined && object.key !== null) {
            message.key = object.key;
        }
        else {
            message.key = "";
        }
        if (object.value !== undefined && object.value !== null) {
            message.value = object.value;
        }
        else {
            message.value = "";
        }
        return message;
    },
};
const baseResponse_ExtensionsEntry = { key: "" };
export const Response_ExtensionsEntry = {
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
        const message = {
            ...baseResponse_ExtensionsEntry,
        };
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
        const message = {
            ...baseResponse_ExtensionsEntry,
        };
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
        const message = {
            ...baseResponse_ExtensionsEntry,
        };
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
const baseInfo = {
    title: "",
    description: "",
    termsOfService: "",
    version: "",
};
export const Info = {
    encode(message, writer = Writer.create()) {
        if (message.title !== "") {
            writer.uint32(10).string(message.title);
        }
        if (message.description !== "") {
            writer.uint32(18).string(message.description);
        }
        if (message.termsOfService !== "") {
            writer.uint32(26).string(message.termsOfService);
        }
        if (message.contact !== undefined) {
            Contact.encode(message.contact, writer.uint32(34).fork()).ldelim();
        }
        if (message.license !== undefined) {
            License.encode(message.license, writer.uint32(42).fork()).ldelim();
        }
        if (message.version !== "") {
            writer.uint32(50).string(message.version);
        }
        Object.entries(message.extensions).forEach(([key, value]) => {
            Info_ExtensionsEntry.encode({ key: key, value }, writer.uint32(58).fork()).ldelim();
        });
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseInfo };
        message.extensions = {};
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.title = reader.string();
                    break;
                case 2:
                    message.description = reader.string();
                    break;
                case 3:
                    message.termsOfService = reader.string();
                    break;
                case 4:
                    message.contact = Contact.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.license = License.decode(reader, reader.uint32());
                    break;
                case 6:
                    message.version = reader.string();
                    break;
                case 7:
                    const entry7 = Info_ExtensionsEntry.decode(reader, reader.uint32());
                    if (entry7.value !== undefined) {
                        message.extensions[entry7.key] = entry7.value;
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
        const message = { ...baseInfo };
        message.extensions = {};
        if (object.title !== undefined && object.title !== null) {
            message.title = String(object.title);
        }
        else {
            message.title = "";
        }
        if (object.description !== undefined && object.description !== null) {
            message.description = String(object.description);
        }
        else {
            message.description = "";
        }
        if (object.termsOfService !== undefined && object.termsOfService !== null) {
            message.termsOfService = String(object.termsOfService);
        }
        else {
            message.termsOfService = "";
        }
        if (object.contact !== undefined && object.contact !== null) {
            message.contact = Contact.fromJSON(object.contact);
        }
        else {
            message.contact = undefined;
        }
        if (object.license !== undefined && object.license !== null) {
            message.license = License.fromJSON(object.license);
        }
        else {
            message.license = undefined;
        }
        if (object.version !== undefined && object.version !== null) {
            message.version = String(object.version);
        }
        else {
            message.version = "";
        }
        if (object.extensions !== undefined && object.extensions !== null) {
            Object.entries(object.extensions).forEach(([key, value]) => {
                message.extensions[key] = Value.fromJSON(value);
            });
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.title !== undefined && (obj.title = message.title);
        message.description !== undefined &&
            (obj.description = message.description);
        message.termsOfService !== undefined &&
            (obj.termsOfService = message.termsOfService);
        message.contact !== undefined &&
            (obj.contact = message.contact
                ? Contact.toJSON(message.contact)
                : undefined);
        message.license !== undefined &&
            (obj.license = message.license
                ? License.toJSON(message.license)
                : undefined);
        message.version !== undefined && (obj.version = message.version);
        obj.extensions = {};
        if (message.extensions) {
            Object.entries(message.extensions).forEach(([k, v]) => {
                obj.extensions[k] = Value.toJSON(v);
            });
        }
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseInfo };
        message.extensions = {};
        if (object.title !== undefined && object.title !== null) {
            message.title = object.title;
        }
        else {
            message.title = "";
        }
        if (object.description !== undefined && object.description !== null) {
            message.description = object.description;
        }
        else {
            message.description = "";
        }
        if (object.termsOfService !== undefined && object.termsOfService !== null) {
            message.termsOfService = object.termsOfService;
        }
        else {
            message.termsOfService = "";
        }
        if (object.contact !== undefined && object.contact !== null) {
            message.contact = Contact.fromPartial(object.contact);
        }
        else {
            message.contact = undefined;
        }
        if (object.license !== undefined && object.license !== null) {
            message.license = License.fromPartial(object.license);
        }
        else {
            message.license = undefined;
        }
        if (object.version !== undefined && object.version !== null) {
            message.version = object.version;
        }
        else {
            message.version = "";
        }
        if (object.extensions !== undefined && object.extensions !== null) {
            Object.entries(object.extensions).forEach(([key, value]) => {
                if (value !== undefined) {
                    message.extensions[key] = Value.fromPartial(value);
                }
            });
        }
        return message;
    },
};
const baseInfo_ExtensionsEntry = { key: "" };
export const Info_ExtensionsEntry = {
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
        const message = { ...baseInfo_ExtensionsEntry };
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
        const message = { ...baseInfo_ExtensionsEntry };
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
        const message = { ...baseInfo_ExtensionsEntry };
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
const baseContact = { name: "", url: "", email: "" };
export const Contact = {
    encode(message, writer = Writer.create()) {
        if (message.name !== "") {
            writer.uint32(10).string(message.name);
        }
        if (message.url !== "") {
            writer.uint32(18).string(message.url);
        }
        if (message.email !== "") {
            writer.uint32(26).string(message.email);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseContact };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.name = reader.string();
                    break;
                case 2:
                    message.url = reader.string();
                    break;
                case 3:
                    message.email = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseContact };
        if (object.name !== undefined && object.name !== null) {
            message.name = String(object.name);
        }
        else {
            message.name = "";
        }
        if (object.url !== undefined && object.url !== null) {
            message.url = String(object.url);
        }
        else {
            message.url = "";
        }
        if (object.email !== undefined && object.email !== null) {
            message.email = String(object.email);
        }
        else {
            message.email = "";
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.name !== undefined && (obj.name = message.name);
        message.url !== undefined && (obj.url = message.url);
        message.email !== undefined && (obj.email = message.email);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseContact };
        if (object.name !== undefined && object.name !== null) {
            message.name = object.name;
        }
        else {
            message.name = "";
        }
        if (object.url !== undefined && object.url !== null) {
            message.url = object.url;
        }
        else {
            message.url = "";
        }
        if (object.email !== undefined && object.email !== null) {
            message.email = object.email;
        }
        else {
            message.email = "";
        }
        return message;
    },
};
const baseLicense = { name: "", url: "" };
export const License = {
    encode(message, writer = Writer.create()) {
        if (message.name !== "") {
            writer.uint32(10).string(message.name);
        }
        if (message.url !== "") {
            writer.uint32(18).string(message.url);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseLicense };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.name = reader.string();
                    break;
                case 2:
                    message.url = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseLicense };
        if (object.name !== undefined && object.name !== null) {
            message.name = String(object.name);
        }
        else {
            message.name = "";
        }
        if (object.url !== undefined && object.url !== null) {
            message.url = String(object.url);
        }
        else {
            message.url = "";
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.name !== undefined && (obj.name = message.name);
        message.url !== undefined && (obj.url = message.url);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseLicense };
        if (object.name !== undefined && object.name !== null) {
            message.name = object.name;
        }
        else {
            message.name = "";
        }
        if (object.url !== undefined && object.url !== null) {
            message.url = object.url;
        }
        else {
            message.url = "";
        }
        return message;
    },
};
const baseExternalDocumentation = { description: "", url: "" };
export const ExternalDocumentation = {
    encode(message, writer = Writer.create()) {
        if (message.description !== "") {
            writer.uint32(10).string(message.description);
        }
        if (message.url !== "") {
            writer.uint32(18).string(message.url);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseExternalDocumentation };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.description = reader.string();
                    break;
                case 2:
                    message.url = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseExternalDocumentation };
        if (object.description !== undefined && object.description !== null) {
            message.description = String(object.description);
        }
        else {
            message.description = "";
        }
        if (object.url !== undefined && object.url !== null) {
            message.url = String(object.url);
        }
        else {
            message.url = "";
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.description !== undefined &&
            (obj.description = message.description);
        message.url !== undefined && (obj.url = message.url);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseExternalDocumentation };
        if (object.description !== undefined && object.description !== null) {
            message.description = object.description;
        }
        else {
            message.description = "";
        }
        if (object.url !== undefined && object.url !== null) {
            message.url = object.url;
        }
        else {
            message.url = "";
        }
        return message;
    },
};
const baseSchema = { discriminator: "", readOnly: false, example: "" };
export const Schema = {
    encode(message, writer = Writer.create()) {
        if (message.jsonSchema !== undefined) {
            JSONSchema.encode(message.jsonSchema, writer.uint32(10).fork()).ldelim();
        }
        if (message.discriminator !== "") {
            writer.uint32(18).string(message.discriminator);
        }
        if (message.readOnly === true) {
            writer.uint32(24).bool(message.readOnly);
        }
        if (message.externalDocs !== undefined) {
            ExternalDocumentation.encode(message.externalDocs, writer.uint32(42).fork()).ldelim();
        }
        if (message.example !== "") {
            writer.uint32(50).string(message.example);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseSchema };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.jsonSchema = JSONSchema.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.discriminator = reader.string();
                    break;
                case 3:
                    message.readOnly = reader.bool();
                    break;
                case 5:
                    message.externalDocs = ExternalDocumentation.decode(reader, reader.uint32());
                    break;
                case 6:
                    message.example = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseSchema };
        if (object.jsonSchema !== undefined && object.jsonSchema !== null) {
            message.jsonSchema = JSONSchema.fromJSON(object.jsonSchema);
        }
        else {
            message.jsonSchema = undefined;
        }
        if (object.discriminator !== undefined && object.discriminator !== null) {
            message.discriminator = String(object.discriminator);
        }
        else {
            message.discriminator = "";
        }
        if (object.readOnly !== undefined && object.readOnly !== null) {
            message.readOnly = Boolean(object.readOnly);
        }
        else {
            message.readOnly = false;
        }
        if (object.externalDocs !== undefined && object.externalDocs !== null) {
            message.externalDocs = ExternalDocumentation.fromJSON(object.externalDocs);
        }
        else {
            message.externalDocs = undefined;
        }
        if (object.example !== undefined && object.example !== null) {
            message.example = String(object.example);
        }
        else {
            message.example = "";
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.jsonSchema !== undefined &&
            (obj.jsonSchema = message.jsonSchema
                ? JSONSchema.toJSON(message.jsonSchema)
                : undefined);
        message.discriminator !== undefined &&
            (obj.discriminator = message.discriminator);
        message.readOnly !== undefined && (obj.readOnly = message.readOnly);
        message.externalDocs !== undefined &&
            (obj.externalDocs = message.externalDocs
                ? ExternalDocumentation.toJSON(message.externalDocs)
                : undefined);
        message.example !== undefined && (obj.example = message.example);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseSchema };
        if (object.jsonSchema !== undefined && object.jsonSchema !== null) {
            message.jsonSchema = JSONSchema.fromPartial(object.jsonSchema);
        }
        else {
            message.jsonSchema = undefined;
        }
        if (object.discriminator !== undefined && object.discriminator !== null) {
            message.discriminator = object.discriminator;
        }
        else {
            message.discriminator = "";
        }
        if (object.readOnly !== undefined && object.readOnly !== null) {
            message.readOnly = object.readOnly;
        }
        else {
            message.readOnly = false;
        }
        if (object.externalDocs !== undefined && object.externalDocs !== null) {
            message.externalDocs = ExternalDocumentation.fromPartial(object.externalDocs);
        }
        else {
            message.externalDocs = undefined;
        }
        if (object.example !== undefined && object.example !== null) {
            message.example = object.example;
        }
        else {
            message.example = "";
        }
        return message;
    },
};
const baseJSONSchema = {
    ref: "",
    title: "",
    description: "",
    default: "",
    readOnly: false,
    example: "",
    multipleOf: 0,
    maximum: 0,
    exclusiveMaximum: false,
    minimum: 0,
    exclusiveMinimum: false,
    maxLength: 0,
    minLength: 0,
    pattern: "",
    maxItems: 0,
    minItems: 0,
    uniqueItems: false,
    maxProperties: 0,
    minProperties: 0,
    required: "",
    array: "",
    type: 0,
    format: "",
    enum: "",
};
export const JSONSchema = {
    encode(message, writer = Writer.create()) {
        if (message.ref !== "") {
            writer.uint32(26).string(message.ref);
        }
        if (message.title !== "") {
            writer.uint32(42).string(message.title);
        }
        if (message.description !== "") {
            writer.uint32(50).string(message.description);
        }
        if (message.default !== "") {
            writer.uint32(58).string(message.default);
        }
        if (message.readOnly === true) {
            writer.uint32(64).bool(message.readOnly);
        }
        if (message.example !== "") {
            writer.uint32(74).string(message.example);
        }
        if (message.multipleOf !== 0) {
            writer.uint32(81).double(message.multipleOf);
        }
        if (message.maximum !== 0) {
            writer.uint32(89).double(message.maximum);
        }
        if (message.exclusiveMaximum === true) {
            writer.uint32(96).bool(message.exclusiveMaximum);
        }
        if (message.minimum !== 0) {
            writer.uint32(105).double(message.minimum);
        }
        if (message.exclusiveMinimum === true) {
            writer.uint32(112).bool(message.exclusiveMinimum);
        }
        if (message.maxLength !== 0) {
            writer.uint32(120).uint64(message.maxLength);
        }
        if (message.minLength !== 0) {
            writer.uint32(128).uint64(message.minLength);
        }
        if (message.pattern !== "") {
            writer.uint32(138).string(message.pattern);
        }
        if (message.maxItems !== 0) {
            writer.uint32(160).uint64(message.maxItems);
        }
        if (message.minItems !== 0) {
            writer.uint32(168).uint64(message.minItems);
        }
        if (message.uniqueItems === true) {
            writer.uint32(176).bool(message.uniqueItems);
        }
        if (message.maxProperties !== 0) {
            writer.uint32(192).uint64(message.maxProperties);
        }
        if (message.minProperties !== 0) {
            writer.uint32(200).uint64(message.minProperties);
        }
        for (const v of message.required) {
            writer.uint32(210).string(v);
        }
        for (const v of message.array) {
            writer.uint32(274).string(v);
        }
        writer.uint32(282).fork();
        for (const v of message.type) {
            writer.int32(v);
        }
        writer.ldelim();
        if (message.format !== "") {
            writer.uint32(290).string(message.format);
        }
        for (const v of message.enum) {
            writer.uint32(370).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseJSONSchema };
        message.required = [];
        message.array = [];
        message.type = [];
        message.enum = [];
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 3:
                    message.ref = reader.string();
                    break;
                case 5:
                    message.title = reader.string();
                    break;
                case 6:
                    message.description = reader.string();
                    break;
                case 7:
                    message.default = reader.string();
                    break;
                case 8:
                    message.readOnly = reader.bool();
                    break;
                case 9:
                    message.example = reader.string();
                    break;
                case 10:
                    message.multipleOf = reader.double();
                    break;
                case 11:
                    message.maximum = reader.double();
                    break;
                case 12:
                    message.exclusiveMaximum = reader.bool();
                    break;
                case 13:
                    message.minimum = reader.double();
                    break;
                case 14:
                    message.exclusiveMinimum = reader.bool();
                    break;
                case 15:
                    message.maxLength = longToNumber(reader.uint64());
                    break;
                case 16:
                    message.minLength = longToNumber(reader.uint64());
                    break;
                case 17:
                    message.pattern = reader.string();
                    break;
                case 20:
                    message.maxItems = longToNumber(reader.uint64());
                    break;
                case 21:
                    message.minItems = longToNumber(reader.uint64());
                    break;
                case 22:
                    message.uniqueItems = reader.bool();
                    break;
                case 24:
                    message.maxProperties = longToNumber(reader.uint64());
                    break;
                case 25:
                    message.minProperties = longToNumber(reader.uint64());
                    break;
                case 26:
                    message.required.push(reader.string());
                    break;
                case 34:
                    message.array.push(reader.string());
                    break;
                case 35:
                    if ((tag & 7) === 2) {
                        const end2 = reader.uint32() + reader.pos;
                        while (reader.pos < end2) {
                            message.type.push(reader.int32());
                        }
                    }
                    else {
                        message.type.push(reader.int32());
                    }
                    break;
                case 36:
                    message.format = reader.string();
                    break;
                case 46:
                    message.enum.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseJSONSchema };
        message.required = [];
        message.array = [];
        message.type = [];
        message.enum = [];
        if (object.ref !== undefined && object.ref !== null) {
            message.ref = String(object.ref);
        }
        else {
            message.ref = "";
        }
        if (object.title !== undefined && object.title !== null) {
            message.title = String(object.title);
        }
        else {
            message.title = "";
        }
        if (object.description !== undefined && object.description !== null) {
            message.description = String(object.description);
        }
        else {
            message.description = "";
        }
        if (object.default !== undefined && object.default !== null) {
            message.default = String(object.default);
        }
        else {
            message.default = "";
        }
        if (object.readOnly !== undefined && object.readOnly !== null) {
            message.readOnly = Boolean(object.readOnly);
        }
        else {
            message.readOnly = false;
        }
        if (object.example !== undefined && object.example !== null) {
            message.example = String(object.example);
        }
        else {
            message.example = "";
        }
        if (object.multipleOf !== undefined && object.multipleOf !== null) {
            message.multipleOf = Number(object.multipleOf);
        }
        else {
            message.multipleOf = 0;
        }
        if (object.maximum !== undefined && object.maximum !== null) {
            message.maximum = Number(object.maximum);
        }
        else {
            message.maximum = 0;
        }
        if (object.exclusiveMaximum !== undefined &&
            object.exclusiveMaximum !== null) {
            message.exclusiveMaximum = Boolean(object.exclusiveMaximum);
        }
        else {
            message.exclusiveMaximum = false;
        }
        if (object.minimum !== undefined && object.minimum !== null) {
            message.minimum = Number(object.minimum);
        }
        else {
            message.minimum = 0;
        }
        if (object.exclusiveMinimum !== undefined &&
            object.exclusiveMinimum !== null) {
            message.exclusiveMinimum = Boolean(object.exclusiveMinimum);
        }
        else {
            message.exclusiveMinimum = false;
        }
        if (object.maxLength !== undefined && object.maxLength !== null) {
            message.maxLength = Number(object.maxLength);
        }
        else {
            message.maxLength = 0;
        }
        if (object.minLength !== undefined && object.minLength !== null) {
            message.minLength = Number(object.minLength);
        }
        else {
            message.minLength = 0;
        }
        if (object.pattern !== undefined && object.pattern !== null) {
            message.pattern = String(object.pattern);
        }
        else {
            message.pattern = "";
        }
        if (object.maxItems !== undefined && object.maxItems !== null) {
            message.maxItems = Number(object.maxItems);
        }
        else {
            message.maxItems = 0;
        }
        if (object.minItems !== undefined && object.minItems !== null) {
            message.minItems = Number(object.minItems);
        }
        else {
            message.minItems = 0;
        }
        if (object.uniqueItems !== undefined && object.uniqueItems !== null) {
            message.uniqueItems = Boolean(object.uniqueItems);
        }
        else {
            message.uniqueItems = false;
        }
        if (object.maxProperties !== undefined && object.maxProperties !== null) {
            message.maxProperties = Number(object.maxProperties);
        }
        else {
            message.maxProperties = 0;
        }
        if (object.minProperties !== undefined && object.minProperties !== null) {
            message.minProperties = Number(object.minProperties);
        }
        else {
            message.minProperties = 0;
        }
        if (object.required !== undefined && object.required !== null) {
            for (const e of object.required) {
                message.required.push(String(e));
            }
        }
        if (object.array !== undefined && object.array !== null) {
            for (const e of object.array) {
                message.array.push(String(e));
            }
        }
        if (object.type !== undefined && object.type !== null) {
            for (const e of object.type) {
                message.type.push(jSONSchema_JSONSchemaSimpleTypesFromJSON(e));
            }
        }
        if (object.format !== undefined && object.format !== null) {
            message.format = String(object.format);
        }
        else {
            message.format = "";
        }
        if (object.enum !== undefined && object.enum !== null) {
            for (const e of object.enum) {
                message.enum.push(String(e));
            }
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.ref !== undefined && (obj.ref = message.ref);
        message.title !== undefined && (obj.title = message.title);
        message.description !== undefined &&
            (obj.description = message.description);
        message.default !== undefined && (obj.default = message.default);
        message.readOnly !== undefined && (obj.readOnly = message.readOnly);
        message.example !== undefined && (obj.example = message.example);
        message.multipleOf !== undefined && (obj.multipleOf = message.multipleOf);
        message.maximum !== undefined && (obj.maximum = message.maximum);
        message.exclusiveMaximum !== undefined &&
            (obj.exclusiveMaximum = message.exclusiveMaximum);
        message.minimum !== undefined && (obj.minimum = message.minimum);
        message.exclusiveMinimum !== undefined &&
            (obj.exclusiveMinimum = message.exclusiveMinimum);
        message.maxLength !== undefined && (obj.maxLength = message.maxLength);
        message.minLength !== undefined && (obj.minLength = message.minLength);
        message.pattern !== undefined && (obj.pattern = message.pattern);
        message.maxItems !== undefined && (obj.maxItems = message.maxItems);
        message.minItems !== undefined && (obj.minItems = message.minItems);
        message.uniqueItems !== undefined &&
            (obj.uniqueItems = message.uniqueItems);
        message.maxProperties !== undefined &&
            (obj.maxProperties = message.maxProperties);
        message.minProperties !== undefined &&
            (obj.minProperties = message.minProperties);
        if (message.required) {
            obj.required = message.required.map((e) => e);
        }
        else {
            obj.required = [];
        }
        if (message.array) {
            obj.array = message.array.map((e) => e);
        }
        else {
            obj.array = [];
        }
        if (message.type) {
            obj.type = message.type.map((e) => jSONSchema_JSONSchemaSimpleTypesToJSON(e));
        }
        else {
            obj.type = [];
        }
        message.format !== undefined && (obj.format = message.format);
        if (message.enum) {
            obj.enum = message.enum.map((e) => e);
        }
        else {
            obj.enum = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseJSONSchema };
        message.required = [];
        message.array = [];
        message.type = [];
        message.enum = [];
        if (object.ref !== undefined && object.ref !== null) {
            message.ref = object.ref;
        }
        else {
            message.ref = "";
        }
        if (object.title !== undefined && object.title !== null) {
            message.title = object.title;
        }
        else {
            message.title = "";
        }
        if (object.description !== undefined && object.description !== null) {
            message.description = object.description;
        }
        else {
            message.description = "";
        }
        if (object.default !== undefined && object.default !== null) {
            message.default = object.default;
        }
        else {
            message.default = "";
        }
        if (object.readOnly !== undefined && object.readOnly !== null) {
            message.readOnly = object.readOnly;
        }
        else {
            message.readOnly = false;
        }
        if (object.example !== undefined && object.example !== null) {
            message.example = object.example;
        }
        else {
            message.example = "";
        }
        if (object.multipleOf !== undefined && object.multipleOf !== null) {
            message.multipleOf = object.multipleOf;
        }
        else {
            message.multipleOf = 0;
        }
        if (object.maximum !== undefined && object.maximum !== null) {
            message.maximum = object.maximum;
        }
        else {
            message.maximum = 0;
        }
        if (object.exclusiveMaximum !== undefined &&
            object.exclusiveMaximum !== null) {
            message.exclusiveMaximum = object.exclusiveMaximum;
        }
        else {
            message.exclusiveMaximum = false;
        }
        if (object.minimum !== undefined && object.minimum !== null) {
            message.minimum = object.minimum;
        }
        else {
            message.minimum = 0;
        }
        if (object.exclusiveMinimum !== undefined &&
            object.exclusiveMinimum !== null) {
            message.exclusiveMinimum = object.exclusiveMinimum;
        }
        else {
            message.exclusiveMinimum = false;
        }
        if (object.maxLength !== undefined && object.maxLength !== null) {
            message.maxLength = object.maxLength;
        }
        else {
            message.maxLength = 0;
        }
        if (object.minLength !== undefined && object.minLength !== null) {
            message.minLength = object.minLength;
        }
        else {
            message.minLength = 0;
        }
        if (object.pattern !== undefined && object.pattern !== null) {
            message.pattern = object.pattern;
        }
        else {
            message.pattern = "";
        }
        if (object.maxItems !== undefined && object.maxItems !== null) {
            message.maxItems = object.maxItems;
        }
        else {
            message.maxItems = 0;
        }
        if (object.minItems !== undefined && object.minItems !== null) {
            message.minItems = object.minItems;
        }
        else {
            message.minItems = 0;
        }
        if (object.uniqueItems !== undefined && object.uniqueItems !== null) {
            message.uniqueItems = object.uniqueItems;
        }
        else {
            message.uniqueItems = false;
        }
        if (object.maxProperties !== undefined && object.maxProperties !== null) {
            message.maxProperties = object.maxProperties;
        }
        else {
            message.maxProperties = 0;
        }
        if (object.minProperties !== undefined && object.minProperties !== null) {
            message.minProperties = object.minProperties;
        }
        else {
            message.minProperties = 0;
        }
        if (object.required !== undefined && object.required !== null) {
            for (const e of object.required) {
                message.required.push(e);
            }
        }
        if (object.array !== undefined && object.array !== null) {
            for (const e of object.array) {
                message.array.push(e);
            }
        }
        if (object.type !== undefined && object.type !== null) {
            for (const e of object.type) {
                message.type.push(e);
            }
        }
        if (object.format !== undefined && object.format !== null) {
            message.format = object.format;
        }
        else {
            message.format = "";
        }
        if (object.enum !== undefined && object.enum !== null) {
            for (const e of object.enum) {
                message.enum.push(e);
            }
        }
        return message;
    },
};
const baseTag = { description: "" };
export const Tag = {
    encode(message, writer = Writer.create()) {
        if (message.description !== "") {
            writer.uint32(18).string(message.description);
        }
        if (message.externalDocs !== undefined) {
            ExternalDocumentation.encode(message.externalDocs, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseTag };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 2:
                    message.description = reader.string();
                    break;
                case 3:
                    message.externalDocs = ExternalDocumentation.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseTag };
        if (object.description !== undefined && object.description !== null) {
            message.description = String(object.description);
        }
        else {
            message.description = "";
        }
        if (object.externalDocs !== undefined && object.externalDocs !== null) {
            message.externalDocs = ExternalDocumentation.fromJSON(object.externalDocs);
        }
        else {
            message.externalDocs = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.description !== undefined &&
            (obj.description = message.description);
        message.externalDocs !== undefined &&
            (obj.externalDocs = message.externalDocs
                ? ExternalDocumentation.toJSON(message.externalDocs)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseTag };
        if (object.description !== undefined && object.description !== null) {
            message.description = object.description;
        }
        else {
            message.description = "";
        }
        if (object.externalDocs !== undefined && object.externalDocs !== null) {
            message.externalDocs = ExternalDocumentation.fromPartial(object.externalDocs);
        }
        else {
            message.externalDocs = undefined;
        }
        return message;
    },
};
const baseSecurityDefinitions = {};
export const SecurityDefinitions = {
    encode(message, writer = Writer.create()) {
        Object.entries(message.security).forEach(([key, value]) => {
            SecurityDefinitions_SecurityEntry.encode({ key: key, value }, writer.uint32(10).fork()).ldelim();
        });
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseSecurityDefinitions };
        message.security = {};
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    const entry1 = SecurityDefinitions_SecurityEntry.decode(reader, reader.uint32());
                    if (entry1.value !== undefined) {
                        message.security[entry1.key] = entry1.value;
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
        const message = { ...baseSecurityDefinitions };
        message.security = {};
        if (object.security !== undefined && object.security !== null) {
            Object.entries(object.security).forEach(([key, value]) => {
                message.security[key] = SecurityScheme.fromJSON(value);
            });
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        obj.security = {};
        if (message.security) {
            Object.entries(message.security).forEach(([k, v]) => {
                obj.security[k] = SecurityScheme.toJSON(v);
            });
        }
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseSecurityDefinitions };
        message.security = {};
        if (object.security !== undefined && object.security !== null) {
            Object.entries(object.security).forEach(([key, value]) => {
                if (value !== undefined) {
                    message.security[key] = SecurityScheme.fromPartial(value);
                }
            });
        }
        return message;
    },
};
const baseSecurityDefinitions_SecurityEntry = { key: "" };
export const SecurityDefinitions_SecurityEntry = {
    encode(message, writer = Writer.create()) {
        if (message.key !== "") {
            writer.uint32(10).string(message.key);
        }
        if (message.value !== undefined) {
            SecurityScheme.encode(message.value, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseSecurityDefinitions_SecurityEntry,
        };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.key = reader.string();
                    break;
                case 2:
                    message.value = SecurityScheme.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseSecurityDefinitions_SecurityEntry,
        };
        if (object.key !== undefined && object.key !== null) {
            message.key = String(object.key);
        }
        else {
            message.key = "";
        }
        if (object.value !== undefined && object.value !== null) {
            message.value = SecurityScheme.fromJSON(object.value);
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
            (obj.value = message.value
                ? SecurityScheme.toJSON(message.value)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseSecurityDefinitions_SecurityEntry,
        };
        if (object.key !== undefined && object.key !== null) {
            message.key = object.key;
        }
        else {
            message.key = "";
        }
        if (object.value !== undefined && object.value !== null) {
            message.value = SecurityScheme.fromPartial(object.value);
        }
        else {
            message.value = undefined;
        }
        return message;
    },
};
const baseSecurityScheme = {
    type: 0,
    description: "",
    name: "",
    in: 0,
    flow: 0,
    authorizationUrl: "",
    tokenUrl: "",
};
export const SecurityScheme = {
    encode(message, writer = Writer.create()) {
        if (message.type !== 0) {
            writer.uint32(8).int32(message.type);
        }
        if (message.description !== "") {
            writer.uint32(18).string(message.description);
        }
        if (message.name !== "") {
            writer.uint32(26).string(message.name);
        }
        if (message.in !== 0) {
            writer.uint32(32).int32(message.in);
        }
        if (message.flow !== 0) {
            writer.uint32(40).int32(message.flow);
        }
        if (message.authorizationUrl !== "") {
            writer.uint32(50).string(message.authorizationUrl);
        }
        if (message.tokenUrl !== "") {
            writer.uint32(58).string(message.tokenUrl);
        }
        if (message.scopes !== undefined) {
            Scopes.encode(message.scopes, writer.uint32(66).fork()).ldelim();
        }
        Object.entries(message.extensions).forEach(([key, value]) => {
            SecurityScheme_ExtensionsEntry.encode({ key: key, value }, writer.uint32(74).fork()).ldelim();
        });
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseSecurityScheme };
        message.extensions = {};
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.type = reader.int32();
                    break;
                case 2:
                    message.description = reader.string();
                    break;
                case 3:
                    message.name = reader.string();
                    break;
                case 4:
                    message.in = reader.int32();
                    break;
                case 5:
                    message.flow = reader.int32();
                    break;
                case 6:
                    message.authorizationUrl = reader.string();
                    break;
                case 7:
                    message.tokenUrl = reader.string();
                    break;
                case 8:
                    message.scopes = Scopes.decode(reader, reader.uint32());
                    break;
                case 9:
                    const entry9 = SecurityScheme_ExtensionsEntry.decode(reader, reader.uint32());
                    if (entry9.value !== undefined) {
                        message.extensions[entry9.key] = entry9.value;
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
        const message = { ...baseSecurityScheme };
        message.extensions = {};
        if (object.type !== undefined && object.type !== null) {
            message.type = securityScheme_TypeFromJSON(object.type);
        }
        else {
            message.type = 0;
        }
        if (object.description !== undefined && object.description !== null) {
            message.description = String(object.description);
        }
        else {
            message.description = "";
        }
        if (object.name !== undefined && object.name !== null) {
            message.name = String(object.name);
        }
        else {
            message.name = "";
        }
        if (object.in !== undefined && object.in !== null) {
            message.in = securityScheme_InFromJSON(object.in);
        }
        else {
            message.in = 0;
        }
        if (object.flow !== undefined && object.flow !== null) {
            message.flow = securityScheme_FlowFromJSON(object.flow);
        }
        else {
            message.flow = 0;
        }
        if (object.authorizationUrl !== undefined &&
            object.authorizationUrl !== null) {
            message.authorizationUrl = String(object.authorizationUrl);
        }
        else {
            message.authorizationUrl = "";
        }
        if (object.tokenUrl !== undefined && object.tokenUrl !== null) {
            message.tokenUrl = String(object.tokenUrl);
        }
        else {
            message.tokenUrl = "";
        }
        if (object.scopes !== undefined && object.scopes !== null) {
            message.scopes = Scopes.fromJSON(object.scopes);
        }
        else {
            message.scopes = undefined;
        }
        if (object.extensions !== undefined && object.extensions !== null) {
            Object.entries(object.extensions).forEach(([key, value]) => {
                message.extensions[key] = Value.fromJSON(value);
            });
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.type !== undefined &&
            (obj.type = securityScheme_TypeToJSON(message.type));
        message.description !== undefined &&
            (obj.description = message.description);
        message.name !== undefined && (obj.name = message.name);
        message.in !== undefined && (obj.in = securityScheme_InToJSON(message.in));
        message.flow !== undefined &&
            (obj.flow = securityScheme_FlowToJSON(message.flow));
        message.authorizationUrl !== undefined &&
            (obj.authorizationUrl = message.authorizationUrl);
        message.tokenUrl !== undefined && (obj.tokenUrl = message.tokenUrl);
        message.scopes !== undefined &&
            (obj.scopes = message.scopes ? Scopes.toJSON(message.scopes) : undefined);
        obj.extensions = {};
        if (message.extensions) {
            Object.entries(message.extensions).forEach(([k, v]) => {
                obj.extensions[k] = Value.toJSON(v);
            });
        }
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseSecurityScheme };
        message.extensions = {};
        if (object.type !== undefined && object.type !== null) {
            message.type = object.type;
        }
        else {
            message.type = 0;
        }
        if (object.description !== undefined && object.description !== null) {
            message.description = object.description;
        }
        else {
            message.description = "";
        }
        if (object.name !== undefined && object.name !== null) {
            message.name = object.name;
        }
        else {
            message.name = "";
        }
        if (object.in !== undefined && object.in !== null) {
            message.in = object.in;
        }
        else {
            message.in = 0;
        }
        if (object.flow !== undefined && object.flow !== null) {
            message.flow = object.flow;
        }
        else {
            message.flow = 0;
        }
        if (object.authorizationUrl !== undefined &&
            object.authorizationUrl !== null) {
            message.authorizationUrl = object.authorizationUrl;
        }
        else {
            message.authorizationUrl = "";
        }
        if (object.tokenUrl !== undefined && object.tokenUrl !== null) {
            message.tokenUrl = object.tokenUrl;
        }
        else {
            message.tokenUrl = "";
        }
        if (object.scopes !== undefined && object.scopes !== null) {
            message.scopes = Scopes.fromPartial(object.scopes);
        }
        else {
            message.scopes = undefined;
        }
        if (object.extensions !== undefined && object.extensions !== null) {
            Object.entries(object.extensions).forEach(([key, value]) => {
                if (value !== undefined) {
                    message.extensions[key] = Value.fromPartial(value);
                }
            });
        }
        return message;
    },
};
const baseSecurityScheme_ExtensionsEntry = { key: "" };
export const SecurityScheme_ExtensionsEntry = {
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
        const message = {
            ...baseSecurityScheme_ExtensionsEntry,
        };
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
        const message = {
            ...baseSecurityScheme_ExtensionsEntry,
        };
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
        const message = {
            ...baseSecurityScheme_ExtensionsEntry,
        };
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
const baseSecurityRequirement = {};
export const SecurityRequirement = {
    encode(message, writer = Writer.create()) {
        Object.entries(message.securityRequirement).forEach(([key, value]) => {
            SecurityRequirement_SecurityRequirementEntry.encode({ key: key, value }, writer.uint32(10).fork()).ldelim();
        });
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseSecurityRequirement };
        message.securityRequirement = {};
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    const entry1 = SecurityRequirement_SecurityRequirementEntry.decode(reader, reader.uint32());
                    if (entry1.value !== undefined) {
                        message.securityRequirement[entry1.key] = entry1.value;
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
        const message = { ...baseSecurityRequirement };
        message.securityRequirement = {};
        if (object.securityRequirement !== undefined &&
            object.securityRequirement !== null) {
            Object.entries(object.securityRequirement).forEach(([key, value]) => {
                message.securityRequirement[key] = SecurityRequirement_SecurityRequirementValue.fromJSON(value);
            });
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        obj.securityRequirement = {};
        if (message.securityRequirement) {
            Object.entries(message.securityRequirement).forEach(([k, v]) => {
                obj.securityRequirement[k] = SecurityRequirement_SecurityRequirementValue.toJSON(v);
            });
        }
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseSecurityRequirement };
        message.securityRequirement = {};
        if (object.securityRequirement !== undefined &&
            object.securityRequirement !== null) {
            Object.entries(object.securityRequirement).forEach(([key, value]) => {
                if (value !== undefined) {
                    message.securityRequirement[key] = SecurityRequirement_SecurityRequirementValue.fromPartial(value);
                }
            });
        }
        return message;
    },
};
const baseSecurityRequirement_SecurityRequirementValue = { scope: "" };
export const SecurityRequirement_SecurityRequirementValue = {
    encode(message, writer = Writer.create()) {
        for (const v of message.scope) {
            writer.uint32(10).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseSecurityRequirement_SecurityRequirementValue,
        };
        message.scope = [];
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.scope.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseSecurityRequirement_SecurityRequirementValue,
        };
        message.scope = [];
        if (object.scope !== undefined && object.scope !== null) {
            for (const e of object.scope) {
                message.scope.push(String(e));
            }
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        if (message.scope) {
            obj.scope = message.scope.map((e) => e);
        }
        else {
            obj.scope = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseSecurityRequirement_SecurityRequirementValue,
        };
        message.scope = [];
        if (object.scope !== undefined && object.scope !== null) {
            for (const e of object.scope) {
                message.scope.push(e);
            }
        }
        return message;
    },
};
const baseSecurityRequirement_SecurityRequirementEntry = { key: "" };
export const SecurityRequirement_SecurityRequirementEntry = {
    encode(message, writer = Writer.create()) {
        if (message.key !== "") {
            writer.uint32(10).string(message.key);
        }
        if (message.value !== undefined) {
            SecurityRequirement_SecurityRequirementValue.encode(message.value, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseSecurityRequirement_SecurityRequirementEntry,
        };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.key = reader.string();
                    break;
                case 2:
                    message.value = SecurityRequirement_SecurityRequirementValue.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseSecurityRequirement_SecurityRequirementEntry,
        };
        if (object.key !== undefined && object.key !== null) {
            message.key = String(object.key);
        }
        else {
            message.key = "";
        }
        if (object.value !== undefined && object.value !== null) {
            message.value = SecurityRequirement_SecurityRequirementValue.fromJSON(object.value);
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
            (obj.value = message.value
                ? SecurityRequirement_SecurityRequirementValue.toJSON(message.value)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseSecurityRequirement_SecurityRequirementEntry,
        };
        if (object.key !== undefined && object.key !== null) {
            message.key = object.key;
        }
        else {
            message.key = "";
        }
        if (object.value !== undefined && object.value !== null) {
            message.value = SecurityRequirement_SecurityRequirementValue.fromPartial(object.value);
        }
        else {
            message.value = undefined;
        }
        return message;
    },
};
const baseScopes = {};
export const Scopes = {
    encode(message, writer = Writer.create()) {
        Object.entries(message.scope).forEach(([key, value]) => {
            Scopes_ScopeEntry.encode({ key: key, value }, writer.uint32(10).fork()).ldelim();
        });
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseScopes };
        message.scope = {};
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    const entry1 = Scopes_ScopeEntry.decode(reader, reader.uint32());
                    if (entry1.value !== undefined) {
                        message.scope[entry1.key] = entry1.value;
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
        const message = { ...baseScopes };
        message.scope = {};
        if (object.scope !== undefined && object.scope !== null) {
            Object.entries(object.scope).forEach(([key, value]) => {
                message.scope[key] = String(value);
            });
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        obj.scope = {};
        if (message.scope) {
            Object.entries(message.scope).forEach(([k, v]) => {
                obj.scope[k] = v;
            });
        }
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseScopes };
        message.scope = {};
        if (object.scope !== undefined && object.scope !== null) {
            Object.entries(object.scope).forEach(([key, value]) => {
                if (value !== undefined) {
                    message.scope[key] = String(value);
                }
            });
        }
        return message;
    },
};
const baseScopes_ScopeEntry = { key: "", value: "" };
export const Scopes_ScopeEntry = {
    encode(message, writer = Writer.create()) {
        if (message.key !== "") {
            writer.uint32(10).string(message.key);
        }
        if (message.value !== "") {
            writer.uint32(18).string(message.value);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseScopes_ScopeEntry };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.key = reader.string();
                    break;
                case 2:
                    message.value = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseScopes_ScopeEntry };
        if (object.key !== undefined && object.key !== null) {
            message.key = String(object.key);
        }
        else {
            message.key = "";
        }
        if (object.value !== undefined && object.value !== null) {
            message.value = String(object.value);
        }
        else {
            message.value = "";
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.key !== undefined && (obj.key = message.key);
        message.value !== undefined && (obj.value = message.value);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseScopes_ScopeEntry };
        if (object.key !== undefined && object.key !== null) {
            message.key = object.key;
        }
        else {
            message.key = "";
        }
        if (object.value !== undefined && object.value !== null) {
            message.value = object.value;
        }
        else {
            message.value = "";
        }
        return message;
    },
};
var globalThis = (() => {
    if (typeof globalThis !== "undefined")
        return globalThis;
    if (typeof self !== "undefined")
        return self;
    if (typeof window !== "undefined")
        return window;
    if (typeof global !== "undefined")
        return global;
    throw "Unable to locate global object";
})();
function longToNumber(long) {
    if (long.gt(Number.MAX_SAFE_INTEGER)) {
        throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
    }
    return long.toNumber();
}
if (util.Long !== Long) {
    util.Long = Long;
    configure();
}
