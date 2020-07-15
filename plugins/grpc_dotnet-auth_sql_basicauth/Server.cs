using System;
using System.Threading.Tasks;
using Grpc.Core;

using System.Data.SqlClient;

class DispatcherImpl : Coprocess.Dispatcher.DispatcherBase {
  public DispatcherImpl()
  {
    Console.WriteLine("Instantiating DispatcherImpl");
    AuthLayer.Init();
  }

  // This is the main method, it will be called for every hook.
  public override Task<Coprocess.Object> Dispatch(Coprocess.Object thisObject, ServerCallContext context)
  {
    Console.WriteLine("Receiving object: " + thisObject.ToString());

    var hook = this.GetType().GetMethod(thisObject.HookName);

    // If DispatcherImpl doesn't implement this hook...
    if(hook == null) {
      Console.WriteLine("Hook name: " + thisObject.HookName + " (not implemented!)");
      // We return the unmodified request object:
      return Task.FromResult(thisObject);
    };

    Console.WriteLine("Hook name: " + thisObject.HookName + " (implemented)");

    // This will dynamically invoke our hook method, and cast the returned object to the required type:
    var output = hook.Invoke(this, new object[]{thisObject, context});
    return (Task<Coprocess.Object>)output;
  }

  // MyPreMiddleware is a sample PRE hook:
  public Task<Coprocess.Object> MyPreMiddleware(Coprocess.Object thisObject, ServerCallContext context)
  {
    Console.WriteLine("Calling MyPreMiddleware.");
    // Injecting a header!
    thisObject.Request.SetHeaders["my-header"] = "my-value";
    return Task.FromResult(thisObject);
  }

  // MyPostMiddleware is a sample POST hook:
  public Task<Coprocess.Object> MyPostMiddleware(Coprocess.Object thisObject, ServerCallContext context)
  {
    Console.WriteLine("Calling MyPostMiddleware.");
    Console.WriteLine(thisObject.Session);
    Console.WriteLine(thisObject);
    return Task.FromResult(thisObject);
  }

  // MyAuthCheck is a sample authentication method, it will interact with AuthLayer (MSSQL database).
  public Task<Coprocess.Object> MyAuthCheck(Coprocess.Object thisObject, ServerCallContext context)
  {
    // Request.Headers contains all the request headers, we retrieve the authorization token:
    var token = thisObject.Request.Headers["Authorization"];
    Console.WriteLine("Calling MyAuthCheck with token = " + token);

    if( !AuthLayer.active ) {
      Console.WriteLine("Rejecting auth! This sample requires a database.");
      return Task.FromResult(thisObject);
    }

    // userData will be null if the authentication is wrong.
    var userData = AuthLayer.Authenticate(token);

    if( userData != null ) {
      Console.WriteLine("Successful auth!");
      var session = new Coprocess.SessionState();
      session.Rate = 1000;
      session.Per = 10;
      session.QuotaMax = 60;
      session.QuotaRenews = 1479033599;
      session.QuotaRemaining = 0;
      session.QuotaRenewalRate = 120;
      session.Expires = 1479033599;

      session.LastUpdated = 1478033599.ToString();

      thisObject.Metadata["token"] = token;
      thisObject.Session = session;
      return Task.FromResult(thisObject);
    }

    Console.WriteLine("Rejecting auth!");
    return Task.FromResult(thisObject);
  }
}
