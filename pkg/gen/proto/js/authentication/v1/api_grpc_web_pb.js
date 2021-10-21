/**
 * @fileoverview gRPC-Web generated client stub for authentication.v1
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.authentication = {};
proto.authentication.v1 = require('./api_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.authentication.v1.AuthenticationServiceClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.authentication.v1.AuthenticationServicePromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.authentication.v1.CreateAccountRequest,
 *   !proto.authentication.v1.CreateAccountResponse>}
 */
const methodDescriptor_AuthenticationService_CreateAccount = new grpc.web.MethodDescriptor(
  '/authentication.v1.AuthenticationService/CreateAccount',
  grpc.web.MethodType.UNARY,
  proto.authentication.v1.CreateAccountRequest,
  proto.authentication.v1.CreateAccountResponse,
  /**
   * @param {!proto.authentication.v1.CreateAccountRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.authentication.v1.CreateAccountResponse.deserializeBinary
);


/**
 * @param {!proto.authentication.v1.CreateAccountRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.authentication.v1.CreateAccountResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.authentication.v1.CreateAccountResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.authentication.v1.AuthenticationServiceClient.prototype.createAccount =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/authentication.v1.AuthenticationService/CreateAccount',
      request,
      metadata || {},
      methodDescriptor_AuthenticationService_CreateAccount,
      callback);
};


/**
 * @param {!proto.authentication.v1.CreateAccountRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.authentication.v1.CreateAccountResponse>}
 *     Promise that resolves to the response
 */
proto.authentication.v1.AuthenticationServicePromiseClient.prototype.createAccount =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/authentication.v1.AuthenticationService/CreateAccount',
      request,
      metadata || {},
      methodDescriptor_AuthenticationService_CreateAccount);
};


module.exports = proto.authentication.v1;

