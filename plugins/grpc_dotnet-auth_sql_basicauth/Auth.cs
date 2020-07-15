using System;
using System.Data;
using System.Data.SqlClient;
using System.Collections.Generic;
using Microsoft.Framework.Configuration;
using Microsoft.Framework.Configuration.Json;

public class AuthLayer {
  private static string connectionString = "";
  private static SqlConnection connection;

  public static bool active = false;

  public static void Init() {
    Console.WriteLine("Initializing AuthLayer.");

    try {
      var builder = new ConfigurationBuilder().AddJsonFile("settings.json");
      var config = builder.Build();
      connectionString = config["connection_string"];
    } catch {
      Console.WriteLine("No settings.json file exists, check settings.sample.json! Auth middleware will be disabled.");
      return;
    }

    try {
      connection = new SqlConnection(connectionString);
      connection.Open();
      active = true;
    } catch {
      Console.WriteLine("Couldn't establish a database connection! Authentication middleware will reject all the requests.");
    }
  }

  public static Dictionary<string, string> Authenticate(string token) {
    var query = "SELECT id, username FROM users WHERE token = @TOKEN";
    SqlCommand cmd = new SqlCommand(query, connection);
    cmd.Parameters.Add("@TOKEN", SqlDbType.VarChar);
    cmd.Parameters["@TOKEN"].Value = token;

    var found = false;
    var userData = new Dictionary<string,string>();

    using (SqlDataReader reader = cmd.ExecuteReader())
    {
      while (reader.Read())
      {
        userData.Add("id", reader.GetInt32(0).ToString());
        userData.Add("username", reader.GetString(1));
        found = true;
      }
    }

    if (!found) {
      return null;
    }

    return userData;
  }

}
