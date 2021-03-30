import { Writer, Reader } from "protobufjs/minimal";
import { Value } from "../../google/protobuf/struct";
export declare const protobufPackage = "grpc.gateway.protoc_gen_openapiv2.options";
/**
 * Scheme describes the schemes supported by the OpenAPI Swagger
 * and Operation objects.
 */
export declare enum Scheme {
    UNKNOWN = 0,
    HTTP = 1,
    HTTPS = 2,
    WS = 3,
    WSS = 4,
    UNRECOGNIZED = -1
}
export declare function schemeFromJSON(object: any): Scheme;
export declare function schemeToJSON(object: Scheme): string;
/**
 * `Swagger` is a representation of OpenAPI v2 specification's Swagger object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#swaggerObject
 *
 * Example:
 *
 *  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
 *    info: {
 *      title: "Echo API";
 *      version: "1.0";
 *      description: ";
 *      contact: {
 *        name: "gRPC-Gateway project";
 *        url: "https://github.com/grpc-ecosystem/grpc-gateway";
 *        email: "none@example.com";
 *      };
 *      license: {
 *        name: "BSD 3-Clause License";
 *        url: "https://github.com/grpc-ecosystem/grpc-gateway/blob/master/LICENSE.txt";
 *      };
 *    };
 *    schemes: HTTPS;
 *    consumes: "application/json";
 *    produces: "application/json";
 *  };
 */
export interface Swagger {
    /**
     * Specifies the OpenAPI Specification version being used. It can be
     * used by the OpenAPI UI and other clients to interpret the API listing. The
     * value MUST be "2.0".
     */
    swagger: string;
    /**
     * Provides metadata about the API. The metadata can be used by the
     * clients if needed.
     */
    info: Info | undefined;
    /**
     * The host (name or ip) serving the API. This MUST be the host only and does
     * not include the scheme nor sub-paths. It MAY include a port. If the host is
     * not included, the host serving the documentation is to be used (including
     * the port). The host does not support path templating.
     */
    host: string;
    /**
     * The base path on which the API is served, which is relative to the host. If
     * it is not included, the API is served directly under the host. The value
     * MUST start with a leading slash (/). The basePath does not support path
     * templating.
     * Note that using `base_path` does not change the endpoint paths that are
     * generated in the resulting OpenAPI file. If you wish to use `base_path`
     * with relatively generated OpenAPI paths, the `base_path` prefix must be
     * manually removed from your `google.api.http` paths and your code changed to
     * serve the API from the `base_path`.
     */
    basePath: string;
    /**
     * The transfer protocol of the API. Values MUST be from the list: "http",
     * "https", "ws", "wss". If the schemes is not included, the default scheme to
     * be used is the one used to access the OpenAPI definition itself.
     */
    schemes: Scheme[];
    /**
     * A list of MIME types the APIs can consume. This is global to all APIs but
     * can be overridden on specific API calls. Value MUST be as described under
     * Mime Types.
     */
    consumes: string[];
    /**
     * A list of MIME types the APIs can produce. This is global to all APIs but
     * can be overridden on specific API calls. Value MUST be as described under
     * Mime Types.
     */
    produces: string[];
    /**
     * An object to hold responses that can be used across operations. This
     * property does not define global responses for all operations.
     */
    responses: {
        [key: string]: Response;
    };
    /** Security scheme definitions that can be used across the specification. */
    securityDefinitions: SecurityDefinitions | undefined;
    /**
     * A declaration of which security schemes are applied for the API as a whole.
     * The list of values describes alternative security schemes that can be used
     * (that is, there is a logical OR between the security requirements).
     * Individual operations can override this definition.
     */
    security: SecurityRequirement[];
    /** Additional external documentation. */
    externalDocs: ExternalDocumentation | undefined;
    extensions: {
        [key: string]: Value;
    };
}
export interface Swagger_ResponsesEntry {
    key: string;
    value: Response | undefined;
}
export interface Swagger_ExtensionsEntry {
    key: string;
    value: Value | undefined;
}
/**
 * `Operation` is a representation of OpenAPI v2 specification's Operation object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#operationObject
 *
 * Example:
 *
 *  service EchoService {
 *    rpc Echo(SimpleMessage) returns (SimpleMessage) {
 *      option (google.api.http) = {
 *        get: "/v1/example/echo/{id}"
 *      };
 *
 *      option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
 *        summary: "Get a message.";
 *        operation_id: "getMessage";
 *        tags: "echo";
 *        responses: {
 *          key: "200"
 *            value: {
 *            description: "OK";
 *          }
 *        }
 *      };
 *    }
 *  }
 */
export interface Operation {
    /**
     * A list of tags for API documentation control. Tags can be used for logical
     * grouping of operations by resources or any other qualifier.
     */
    tags: string[];
    /**
     * A short summary of what the operation does. For maximum readability in the
     * swagger-ui, this field SHOULD be less than 120 characters.
     */
    summary: string;
    /**
     * A verbose explanation of the operation behavior. GFM syntax can be used for
     * rich text representation.
     */
    description: string;
    /** Additional external documentation for this operation. */
    externalDocs: ExternalDocumentation | undefined;
    /**
     * Unique string used to identify the operation. The id MUST be unique among
     * all operations described in the API. Tools and libraries MAY use the
     * operationId to uniquely identify an operation, therefore, it is recommended
     * to follow common programming naming conventions.
     */
    operationId: string;
    /**
     * A list of MIME types the operation can consume. This overrides the consumes
     * definition at the OpenAPI Object. An empty value MAY be used to clear the
     * global definition. Value MUST be as described under Mime Types.
     */
    consumes: string[];
    /**
     * A list of MIME types the operation can produce. This overrides the produces
     * definition at the OpenAPI Object. An empty value MAY be used to clear the
     * global definition. Value MUST be as described under Mime Types.
     */
    produces: string[];
    /**
     * The list of possible responses as they are returned from executing this
     * operation.
     */
    responses: {
        [key: string]: Response;
    };
    /**
     * The transfer protocol for the operation. Values MUST be from the list:
     * "http", "https", "ws", "wss". The value overrides the OpenAPI Object
     * schemes definition.
     */
    schemes: Scheme[];
    /**
     * Declares this operation to be deprecated. Usage of the declared operation
     * should be refrained. Default value is false.
     */
    deprecated: boolean;
    /**
     * A declaration of which security schemes are applied for this operation. The
     * list of values describes alternative security schemes that can be used
     * (that is, there is a logical OR between the security requirements). This
     * definition overrides any declared top-level security. To remove a top-level
     * security declaration, an empty array can be used.
     */
    security: SecurityRequirement[];
    extensions: {
        [key: string]: Value;
    };
}
export interface Operation_ResponsesEntry {
    key: string;
    value: Response | undefined;
}
export interface Operation_ExtensionsEntry {
    key: string;
    value: Value | undefined;
}
/**
 * `Header` is a representation of OpenAPI v2 specification's Header object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#headerObject
 */
export interface Header {
    /** `Description` is a short description of the header. */
    description: string;
    /** The type of the object. The value MUST be one of "string", "number", "integer", or "boolean". The "array" type is not supported. */
    type: string;
    /** `Format` The extending format for the previously mentioned type. */
    format: string;
    /**
     * `Default` Declares the value of the header that the server will use if none is provided.
     * See: https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-6.2.
     * Unlike JSON Schema this value MUST conform to the defined type for the header.
     */
    default: string;
    /** 'Pattern' See https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-5.2.3. */
    pattern: string;
}
/**
 * `Response` is a representation of OpenAPI v2 specification's Response object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#responseObject
 */
export interface Response {
    /**
     * `Description` is a short description of the response.
     * GFM syntax can be used for rich text representation.
     */
    description: string;
    /**
     * `Schema` optionally defines the structure of the response.
     * If `Schema` is not provided, it means there is no content to the response.
     */
    schema: Schema | undefined;
    /**
     * `Headers` A list of headers that are sent with the response.
     * `Header` name is expected to be a string in the canonical format of the MIME header key
     * See: https://golang.org/pkg/net/textproto/#CanonicalMIMEHeaderKey
     */
    headers: {
        [key: string]: Header;
    };
    /**
     * `Examples` gives per-mimetype response examples.
     * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#example-object
     */
    examples: {
        [key: string]: string;
    };
    extensions: {
        [key: string]: Value;
    };
}
export interface Response_HeadersEntry {
    key: string;
    value: Header | undefined;
}
export interface Response_ExamplesEntry {
    key: string;
    value: string;
}
export interface Response_ExtensionsEntry {
    key: string;
    value: Value | undefined;
}
/**
 * `Info` is a representation of OpenAPI v2 specification's Info object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#infoObject
 *
 * Example:
 *
 *  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
 *    info: {
 *      title: "Echo API";
 *      version: "1.0";
 *      description: ";
 *      contact: {
 *        name: "gRPC-Gateway project";
 *        url: "https://github.com/grpc-ecosystem/grpc-gateway";
 *        email: "none@example.com";
 *      };
 *      license: {
 *        name: "BSD 3-Clause License";
 *        url: "https://github.com/grpc-ecosystem/grpc-gateway/blob/master/LICENSE.txt";
 *      };
 *    };
 *    ...
 *  };
 */
export interface Info {
    /** The title of the application. */
    title: string;
    /**
     * A short description of the application. GFM syntax can be used for rich
     * text representation.
     */
    description: string;
    /** The Terms of Service for the API. */
    termsOfService: string;
    /** The contact information for the exposed API. */
    contact: Contact | undefined;
    /** The license information for the exposed API. */
    license: License | undefined;
    /**
     * Provides the version of the application API (not to be confused
     * with the specification version).
     */
    version: string;
    extensions: {
        [key: string]: Value;
    };
}
export interface Info_ExtensionsEntry {
    key: string;
    value: Value | undefined;
}
/**
 * `Contact` is a representation of OpenAPI v2 specification's Contact object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#contactObject
 *
 * Example:
 *
 *  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
 *    info: {
 *      ...
 *      contact: {
 *        name: "gRPC-Gateway project";
 *        url: "https://github.com/grpc-ecosystem/grpc-gateway";
 *        email: "none@example.com";
 *      };
 *      ...
 *    };
 *    ...
 *  };
 */
export interface Contact {
    /** The identifying name of the contact person/organization. */
    name: string;
    /**
     * The URL pointing to the contact information. MUST be in the format of a
     * URL.
     */
    url: string;
    /**
     * The email address of the contact person/organization. MUST be in the format
     * of an email address.
     */
    email: string;
}
/**
 * `License` is a representation of OpenAPI v2 specification's License object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#licenseObject
 *
 * Example:
 *
 *  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
 *    info: {
 *      ...
 *      license: {
 *        name: "BSD 3-Clause License";
 *        url: "https://github.com/grpc-ecosystem/grpc-gateway/blob/master/LICENSE.txt";
 *      };
 *      ...
 *    };
 *    ...
 *  };
 */
export interface License {
    /** The license name used for the API. */
    name: string;
    /** A URL to the license used for the API. MUST be in the format of a URL. */
    url: string;
}
/**
 * `ExternalDocumentation` is a representation of OpenAPI v2 specification's
 * ExternalDocumentation object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#externalDocumentationObject
 *
 * Example:
 *
 *  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
 *    ...
 *    external_docs: {
 *      description: "More about gRPC-Gateway";
 *      url: "https://github.com/grpc-ecosystem/grpc-gateway";
 *    }
 *    ...
 *  };
 */
export interface ExternalDocumentation {
    /**
     * A short description of the target documentation. GFM syntax can be used for
     * rich text representation.
     */
    description: string;
    /**
     * The URL for the target documentation. Value MUST be in the format
     * of a URL.
     */
    url: string;
}
/**
 * `Schema` is a representation of OpenAPI v2 specification's Schema object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#schemaObject
 */
export interface Schema {
    jsonSchema: JSONSchema | undefined;
    /**
     * Adds support for polymorphism. The discriminator is the schema property
     * name that is used to differentiate between other schema that inherit this
     * schema. The property name used MUST be defined at this schema and it MUST
     * be in the required property list. When used, the value MUST be the name of
     * this schema or any schema that inherits it.
     */
    discriminator: string;
    /**
     * Relevant only for Schema "properties" definitions. Declares the property as
     * "read only". This means that it MAY be sent as part of a response but MUST
     * NOT be sent as part of the request. Properties marked as readOnly being
     * true SHOULD NOT be in the required list of the defined schema. Default
     * value is false.
     */
    readOnly: boolean;
    /** Additional external documentation for this schema. */
    externalDocs: ExternalDocumentation | undefined;
    /**
     * A free-form property to include an example of an instance for this schema in JSON.
     * This is copied verbatim to the output.
     */
    example: string;
}
/**
 * `JSONSchema` represents properties from JSON Schema taken, and as used, in
 * the OpenAPI v2 spec.
 *
 * This includes changes made by OpenAPI v2.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#schemaObject
 *
 * See also: https://cswr.github.io/JsonSchema/spec/basic_types/,
 * https://github.com/json-schema-org/json-schema-spec/blob/master/schema.json
 *
 * Example:
 *
 *  message SimpleMessage {
 *    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
 *      json_schema: {
 *        title: "SimpleMessage"
 *        description: "A simple message."
 *        required: ["id"]
 *      }
 *    };
 *
 *    // Id represents the message identifier.
 *    string id = 1; [
 *        (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = {
 *          {description: "The unique identifier of the simple message."
 *        }];
 *  }
 */
export interface JSONSchema {
    /**
     * Ref is used to define an external reference to include in the message.
     * This could be a fully qualified proto message reference, and that type must
     * be imported into the protofile. If no message is identified, the Ref will
     * be used verbatim in the output.
     * For example:
     *  `ref: ".google.protobuf.Timestamp"`.
     */
    ref: string;
    /** The title of the schema. */
    title: string;
    /** A short description of the schema. */
    description: string;
    default: string;
    readOnly: boolean;
    /**
     * A free-form property to include a JSON example of this field. This is copied
     * verbatim to the output swagger.json. Quotes must be escaped.
     * This property is the same for 2.0 and 3.0.0 https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/3.0.0.md#schemaObject  https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#schemaObject
     */
    example: string;
    multipleOf: number;
    /**
     * Maximum represents an inclusive upper limit for a numeric instance. The
     * value of MUST be a number,
     */
    maximum: number;
    exclusiveMaximum: boolean;
    /**
     * minimum represents an inclusive lower limit for a numeric instance. The
     * value of MUST be a number,
     */
    minimum: number;
    exclusiveMinimum: boolean;
    maxLength: number;
    minLength: number;
    pattern: string;
    maxItems: number;
    minItems: number;
    uniqueItems: boolean;
    maxProperties: number;
    minProperties: number;
    required: string[];
    /** Items in 'array' must be unique. */
    array: string[];
    type: JSONSchema_JSONSchemaSimpleTypes[];
    /** `Format` */
    format: string;
    /** Items in `enum` must be unique https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-5.5.1 */
    enum: string[];
}
export declare enum JSONSchema_JSONSchemaSimpleTypes {
    UNKNOWN = 0,
    ARRAY = 1,
    BOOLEAN = 2,
    INTEGER = 3,
    NULL = 4,
    NUMBER = 5,
    OBJECT = 6,
    STRING = 7,
    UNRECOGNIZED = -1
}
export declare function jSONSchema_JSONSchemaSimpleTypesFromJSON(object: any): JSONSchema_JSONSchemaSimpleTypes;
export declare function jSONSchema_JSONSchemaSimpleTypesToJSON(object: JSONSchema_JSONSchemaSimpleTypes): string;
/**
 * `Tag` is a representation of OpenAPI v2 specification's Tag object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#tagObject
 */
export interface Tag {
    /**
     * A short description for the tag. GFM syntax can be used for rich text
     * representation.
     */
    description: string;
    /** Additional external documentation for this tag. */
    externalDocs: ExternalDocumentation | undefined;
}
/**
 * `SecurityDefinitions` is a representation of OpenAPI v2 specification's
 * Security Definitions object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#securityDefinitionsObject
 *
 * A declaration of the security schemes available to be used in the
 * specification. This does not enforce the security schemes on the operations
 * and only serves to provide the relevant details for each scheme.
 */
export interface SecurityDefinitions {
    /**
     * A single security scheme definition, mapping a "name" to the scheme it
     * defines.
     */
    security: {
        [key: string]: SecurityScheme;
    };
}
export interface SecurityDefinitions_SecurityEntry {
    key: string;
    value: SecurityScheme | undefined;
}
/**
 * `SecurityScheme` is a representation of OpenAPI v2 specification's
 * Security Scheme object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#securitySchemeObject
 *
 * Allows the definition of a security scheme that can be used by the
 * operations. Supported schemes are basic authentication, an API key (either as
 * a header or as a query parameter) and OAuth2's common flows (implicit,
 * password, application and access code).
 */
export interface SecurityScheme {
    /**
     * The type of the security scheme. Valid values are "basic",
     * "apiKey" or "oauth2".
     */
    type: SecurityScheme_Type;
    /** A short description for security scheme. */
    description: string;
    /**
     * The name of the header or query parameter to be used.
     * Valid for apiKey.
     */
    name: string;
    /**
     * The location of the API key. Valid values are "query" or
     * "header".
     * Valid for apiKey.
     */
    in: SecurityScheme_In;
    /**
     * The flow used by the OAuth2 security scheme. Valid values are
     * "implicit", "password", "application" or "accessCode".
     * Valid for oauth2.
     */
    flow: SecurityScheme_Flow;
    /**
     * The authorization URL to be used for this flow. This SHOULD be in
     * the form of a URL.
     * Valid for oauth2/implicit and oauth2/accessCode.
     */
    authorizationUrl: string;
    /**
     * The token URL to be used for this flow. This SHOULD be in the
     * form of a URL.
     * Valid for oauth2/password, oauth2/application and oauth2/accessCode.
     */
    tokenUrl: string;
    /**
     * The available scopes for the OAuth2 security scheme.
     * Valid for oauth2.
     */
    scopes: Scopes | undefined;
    extensions: {
        [key: string]: Value;
    };
}
/**
 * The type of the security scheme. Valid values are "basic",
 * "apiKey" or "oauth2".
 */
export declare enum SecurityScheme_Type {
    TYPE_INVALID = 0,
    TYPE_BASIC = 1,
    TYPE_API_KEY = 2,
    TYPE_OAUTH2 = 3,
    UNRECOGNIZED = -1
}
export declare function securityScheme_TypeFromJSON(object: any): SecurityScheme_Type;
export declare function securityScheme_TypeToJSON(object: SecurityScheme_Type): string;
/** The location of the API key. Valid values are "query" or "header". */
export declare enum SecurityScheme_In {
    IN_INVALID = 0,
    IN_QUERY = 1,
    IN_HEADER = 2,
    UNRECOGNIZED = -1
}
export declare function securityScheme_InFromJSON(object: any): SecurityScheme_In;
export declare function securityScheme_InToJSON(object: SecurityScheme_In): string;
/**
 * The flow used by the OAuth2 security scheme. Valid values are
 * "implicit", "password", "application" or "accessCode".
 */
export declare enum SecurityScheme_Flow {
    FLOW_INVALID = 0,
    FLOW_IMPLICIT = 1,
    FLOW_PASSWORD = 2,
    FLOW_APPLICATION = 3,
    FLOW_ACCESS_CODE = 4,
    UNRECOGNIZED = -1
}
export declare function securityScheme_FlowFromJSON(object: any): SecurityScheme_Flow;
export declare function securityScheme_FlowToJSON(object: SecurityScheme_Flow): string;
export interface SecurityScheme_ExtensionsEntry {
    key: string;
    value: Value | undefined;
}
/**
 * `SecurityRequirement` is a representation of OpenAPI v2 specification's
 * Security Requirement object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#securityRequirementObject
 *
 * Lists the required security schemes to execute this operation. The object can
 * have multiple security schemes declared in it which are all required (that
 * is, there is a logical AND between the schemes).
 *
 * The name used for each property MUST correspond to a security scheme
 * declared in the Security Definitions.
 */
export interface SecurityRequirement {
    /**
     * Each name must correspond to a security scheme which is declared in
     * the Security Definitions. If the security scheme is of type "oauth2",
     * then the value is a list of scope names required for the execution.
     * For other security scheme types, the array MUST be empty.
     */
    securityRequirement: {
        [key: string]: SecurityRequirement_SecurityRequirementValue;
    };
}
/**
 * If the security scheme is of type "oauth2", then the value is a list of
 * scope names required for the execution. For other security scheme types,
 * the array MUST be empty.
 */
export interface SecurityRequirement_SecurityRequirementValue {
    scope: string[];
}
export interface SecurityRequirement_SecurityRequirementEntry {
    key: string;
    value: SecurityRequirement_SecurityRequirementValue | undefined;
}
/**
 * `Scopes` is a representation of OpenAPI v2 specification's Scopes object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#scopesObject
 *
 * Lists the available scopes for an OAuth2 security scheme.
 */
export interface Scopes {
    /**
     * Maps between a name of a scope to a short description of it (as the value
     * of the property).
     */
    scope: {
        [key: string]: string;
    };
}
export interface Scopes_ScopeEntry {
    key: string;
    value: string;
}
export declare const Swagger: {
    encode(message: Swagger, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Swagger;
    fromJSON(object: any): Swagger;
    toJSON(message: Swagger): unknown;
    fromPartial(object: DeepPartial<Swagger>): Swagger;
};
export declare const Swagger_ResponsesEntry: {
    encode(message: Swagger_ResponsesEntry, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Swagger_ResponsesEntry;
    fromJSON(object: any): Swagger_ResponsesEntry;
    toJSON(message: Swagger_ResponsesEntry): unknown;
    fromPartial(object: DeepPartial<Swagger_ResponsesEntry>): Swagger_ResponsesEntry;
};
export declare const Swagger_ExtensionsEntry: {
    encode(message: Swagger_ExtensionsEntry, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Swagger_ExtensionsEntry;
    fromJSON(object: any): Swagger_ExtensionsEntry;
    toJSON(message: Swagger_ExtensionsEntry): unknown;
    fromPartial(object: DeepPartial<Swagger_ExtensionsEntry>): Swagger_ExtensionsEntry;
};
export declare const Operation: {
    encode(message: Operation, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Operation;
    fromJSON(object: any): Operation;
    toJSON(message: Operation): unknown;
    fromPartial(object: DeepPartial<Operation>): Operation;
};
export declare const Operation_ResponsesEntry: {
    encode(message: Operation_ResponsesEntry, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Operation_ResponsesEntry;
    fromJSON(object: any): Operation_ResponsesEntry;
    toJSON(message: Operation_ResponsesEntry): unknown;
    fromPartial(object: DeepPartial<Operation_ResponsesEntry>): Operation_ResponsesEntry;
};
export declare const Operation_ExtensionsEntry: {
    encode(message: Operation_ExtensionsEntry, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Operation_ExtensionsEntry;
    fromJSON(object: any): Operation_ExtensionsEntry;
    toJSON(message: Operation_ExtensionsEntry): unknown;
    fromPartial(object: DeepPartial<Operation_ExtensionsEntry>): Operation_ExtensionsEntry;
};
export declare const Header: {
    encode(message: Header, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Header;
    fromJSON(object: any): Header;
    toJSON(message: Header): unknown;
    fromPartial(object: DeepPartial<Header>): Header;
};
export declare const Response: {
    encode(message: Response, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Response;
    fromJSON(object: any): Response;
    toJSON(message: Response): unknown;
    fromPartial(object: DeepPartial<Response>): Response;
};
export declare const Response_HeadersEntry: {
    encode(message: Response_HeadersEntry, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Response_HeadersEntry;
    fromJSON(object: any): Response_HeadersEntry;
    toJSON(message: Response_HeadersEntry): unknown;
    fromPartial(object: DeepPartial<Response_HeadersEntry>): Response_HeadersEntry;
};
export declare const Response_ExamplesEntry: {
    encode(message: Response_ExamplesEntry, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Response_ExamplesEntry;
    fromJSON(object: any): Response_ExamplesEntry;
    toJSON(message: Response_ExamplesEntry): unknown;
    fromPartial(object: DeepPartial<Response_ExamplesEntry>): Response_ExamplesEntry;
};
export declare const Response_ExtensionsEntry: {
    encode(message: Response_ExtensionsEntry, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Response_ExtensionsEntry;
    fromJSON(object: any): Response_ExtensionsEntry;
    toJSON(message: Response_ExtensionsEntry): unknown;
    fromPartial(object: DeepPartial<Response_ExtensionsEntry>): Response_ExtensionsEntry;
};
export declare const Info: {
    encode(message: Info, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Info;
    fromJSON(object: any): Info;
    toJSON(message: Info): unknown;
    fromPartial(object: DeepPartial<Info>): Info;
};
export declare const Info_ExtensionsEntry: {
    encode(message: Info_ExtensionsEntry, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Info_ExtensionsEntry;
    fromJSON(object: any): Info_ExtensionsEntry;
    toJSON(message: Info_ExtensionsEntry): unknown;
    fromPartial(object: DeepPartial<Info_ExtensionsEntry>): Info_ExtensionsEntry;
};
export declare const Contact: {
    encode(message: Contact, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Contact;
    fromJSON(object: any): Contact;
    toJSON(message: Contact): unknown;
    fromPartial(object: DeepPartial<Contact>): Contact;
};
export declare const License: {
    encode(message: License, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): License;
    fromJSON(object: any): License;
    toJSON(message: License): unknown;
    fromPartial(object: DeepPartial<License>): License;
};
export declare const ExternalDocumentation: {
    encode(message: ExternalDocumentation, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): ExternalDocumentation;
    fromJSON(object: any): ExternalDocumentation;
    toJSON(message: ExternalDocumentation): unknown;
    fromPartial(object: DeepPartial<ExternalDocumentation>): ExternalDocumentation;
};
export declare const Schema: {
    encode(message: Schema, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Schema;
    fromJSON(object: any): Schema;
    toJSON(message: Schema): unknown;
    fromPartial(object: DeepPartial<Schema>): Schema;
};
export declare const JSONSchema: {
    encode(message: JSONSchema, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): JSONSchema;
    fromJSON(object: any): JSONSchema;
    toJSON(message: JSONSchema): unknown;
    fromPartial(object: DeepPartial<JSONSchema>): JSONSchema;
};
export declare const Tag: {
    encode(message: Tag, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Tag;
    fromJSON(object: any): Tag;
    toJSON(message: Tag): unknown;
    fromPartial(object: DeepPartial<Tag>): Tag;
};
export declare const SecurityDefinitions: {
    encode(message: SecurityDefinitions, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): SecurityDefinitions;
    fromJSON(object: any): SecurityDefinitions;
    toJSON(message: SecurityDefinitions): unknown;
    fromPartial(object: DeepPartial<SecurityDefinitions>): SecurityDefinitions;
};
export declare const SecurityDefinitions_SecurityEntry: {
    encode(message: SecurityDefinitions_SecurityEntry, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): SecurityDefinitions_SecurityEntry;
    fromJSON(object: any): SecurityDefinitions_SecurityEntry;
    toJSON(message: SecurityDefinitions_SecurityEntry): unknown;
    fromPartial(object: DeepPartial<SecurityDefinitions_SecurityEntry>): SecurityDefinitions_SecurityEntry;
};
export declare const SecurityScheme: {
    encode(message: SecurityScheme, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): SecurityScheme;
    fromJSON(object: any): SecurityScheme;
    toJSON(message: SecurityScheme): unknown;
    fromPartial(object: DeepPartial<SecurityScheme>): SecurityScheme;
};
export declare const SecurityScheme_ExtensionsEntry: {
    encode(message: SecurityScheme_ExtensionsEntry, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): SecurityScheme_ExtensionsEntry;
    fromJSON(object: any): SecurityScheme_ExtensionsEntry;
    toJSON(message: SecurityScheme_ExtensionsEntry): unknown;
    fromPartial(object: DeepPartial<SecurityScheme_ExtensionsEntry>): SecurityScheme_ExtensionsEntry;
};
export declare const SecurityRequirement: {
    encode(message: SecurityRequirement, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): SecurityRequirement;
    fromJSON(object: any): SecurityRequirement;
    toJSON(message: SecurityRequirement): unknown;
    fromPartial(object: DeepPartial<SecurityRequirement>): SecurityRequirement;
};
export declare const SecurityRequirement_SecurityRequirementValue: {
    encode(message: SecurityRequirement_SecurityRequirementValue, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): SecurityRequirement_SecurityRequirementValue;
    fromJSON(object: any): SecurityRequirement_SecurityRequirementValue;
    toJSON(message: SecurityRequirement_SecurityRequirementValue): unknown;
    fromPartial(object: DeepPartial<SecurityRequirement_SecurityRequirementValue>): SecurityRequirement_SecurityRequirementValue;
};
export declare const SecurityRequirement_SecurityRequirementEntry: {
    encode(message: SecurityRequirement_SecurityRequirementEntry, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): SecurityRequirement_SecurityRequirementEntry;
    fromJSON(object: any): SecurityRequirement_SecurityRequirementEntry;
    toJSON(message: SecurityRequirement_SecurityRequirementEntry): unknown;
    fromPartial(object: DeepPartial<SecurityRequirement_SecurityRequirementEntry>): SecurityRequirement_SecurityRequirementEntry;
};
export declare const Scopes: {
    encode(message: Scopes, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Scopes;
    fromJSON(object: any): Scopes;
    toJSON(message: Scopes): unknown;
    fromPartial(object: DeepPartial<Scopes>): Scopes;
};
export declare const Scopes_ScopeEntry: {
    encode(message: Scopes_ScopeEntry, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Scopes_ScopeEntry;
    fromJSON(object: any): Scopes_ScopeEntry;
    toJSON(message: Scopes_ScopeEntry): unknown;
    fromPartial(object: DeepPartial<Scopes_ScopeEntry>): Scopes_ScopeEntry;
};
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
