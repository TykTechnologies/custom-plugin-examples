package com.tyktechnologies.tykmiddleware;

import com.auth0.jwt.JWT;
import com.auth0.jwt.algorithms.Algorithm;
import com.auth0.jwt.interfaces.DecodedJWT;
import coprocess.CoprocessObject;
import coprocess.CoprocessSessionState;
import coprocess.DispatcherGrpc;

import java.util.Date;
import java.util.HashMap;

public class TykDispatcher extends DispatcherGrpc.DispatcherImplBase {

    @Override
    public void dispatch(CoprocessObject.Object request,
            io.grpc.stub.StreamObserver<CoprocessObject.Object> responseObserver) {

        System.out.println("*** Incoming Request ***");
        System.out.println("Hook name: " + request.getHookName());
        System.out.println("existing metadata: " + request.getMetadataMap());

        final CoprocessObject.Object modifiedRequest = MetaDataInject(request);
        
        responseObserver.onNext(modifiedRequest);
        responseObserver.onCompleted();
    }

    CoprocessObject.Object MetaDataInject(CoprocessObject.Object request) {
        CoprocessObject.Object.Builder builder = request.toBuilder();
        builder.getRequestBuilder().putSetHeaders("submitterId", request.getMetadataMap().get("tyk_key_request_fields").split("\"")[3]);
        return builder.build();
    }
}