using System;
using Grpc.Core;

namespace ConsoleApplication
{
  public class Program
  {
    // Make sure your Tyk instance is pointing to this port (tcp://127.0.0.1:5555):
    const int Port = 5555;

    public static void Main(string[] args)
    {
      Server server = new Server{
        Services = { Coprocess.Dispatcher.BindService(new DispatcherImpl()) },
        Ports = { new ServerPort("localhost", Port, ServerCredentials.Insecure) }
      };
      server.Start();

      Console.WriteLine("gRPC / CP server listening on " + Port);
      Console.WriteLine("Press any key to stop the server...");
      Console.ReadKey();

      server.ShutdownAsync().Wait();
    }
  }
}
