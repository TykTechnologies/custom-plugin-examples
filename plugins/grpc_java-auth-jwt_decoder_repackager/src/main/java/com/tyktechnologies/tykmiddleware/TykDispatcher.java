package com.tyktechnologies.tykmiddleware;

import com.auth0.jwt.JWT;
import com.auth0.jwt.algorithms.Algorithm;
import com.auth0.jwt.interfaces.DecodedJWT;
import coprocess.DispatcherGrpc;
import coprocess.CoprocessObject;
import coprocess.CoprocessSessionState;
import coprocess.CoprocessReturnOverrides;

import java.util.Date;

public class TykDispatcher extends DispatcherGrpc.DispatcherImplBase {

    final String FOOBAR = "foobar";

    @Override
    public void dispatch(CoprocessObject.Object request,
            io.grpc.stub.StreamObserver<CoprocessObject.Object> responseObserver) {

        System.out.println("*** Incoming Request ***");
        System.out.println("Hook name: " + request.getHookName());
        System.out.println("existing metadata: " + request.getMetadataMap());

        final CoprocessObject.Object modifiedRequest = MyAuthHook(request);
        System.out.println("updated metadata: " + modifiedRequest.getMetadataMap());
        responseObserver.onNext(modifiedRequest);

        System.out.println("*** Transformed Request ***");

        responseObserver.onCompleted();
    }

    CoprocessObject.Object MyAuthHook(CoprocessObject.Object request) {
        // print the JWT or THROW
        final String initialJWT = request.getRequest().getHeadersOrThrow("Authorization");
        System.out.println(initialJWT);

        // Decode Existing JWT
        Algorithm algorithmHS = Algorithm.HMAC256("tyk123");
        DecodedJWT decodedJWT = JWT.decode(initialJWT);

        // Create new JWT
        String newJwt = JWT.create()
                .withClaim("role", decodedJWT.getClaim("role").asString())  // Injects "role" from old JWT
                .withClaim("my-new-claim", "value of new claim")                              // Injects new claim
                .withSubject(decodedJWT.getSubject())
                .withIssuedAt(decodedJWT.getIssuedAt())
                .sign(algorithmHS);

        // Set the Tyk Session
        final long expiryTime = (System.currentTimeMillis() / 1000) + 5;
        CoprocessSessionState.SessionState session = CoprocessSessionState.SessionState.newBuilder()
        .setRate(1000.0)
        .setPer(1.0)
        .setIdExtractorDeadline(expiryTime)
        .build();

        CoprocessObject.Object.Builder builder = request.toBuilder();
        builder.getRequestBuilder().putSetHeaders("Authorization", newJwt);
        builder.putMetadata("token", newJwt);
        builder.setSession(session);

        return builder.build();
    }
}
