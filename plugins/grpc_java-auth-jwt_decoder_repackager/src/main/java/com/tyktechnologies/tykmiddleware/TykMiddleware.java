package com.tyktechnologies.tykmiddleware;

import coprocess.DispatcherGrpc;

import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.stub.StreamObserver;
import java.io.IOException;
import java.util.logging.Level;
import java.util.logging.Logger;

public class TykMiddleware {

    private static final Logger logger = Logger.getLogger(TykMiddleware.class.getName());
    static Server server;
    static int port = 5555;

    public static void main(String[] args) throws IOException, InterruptedException {
        System.out.println("Initializing gRPC server.");
        
        server = ServerBuilder.forPort(port)
                .addService(new TykDispatcher())
                .build()
                .start();
        
        blockUntilShutdown();

    }

    static void blockUntilShutdown() throws InterruptedException {
        if (server != null) {
            server.awaitTermination();
        }
    }
}
