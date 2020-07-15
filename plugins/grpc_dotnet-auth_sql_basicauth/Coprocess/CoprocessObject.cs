// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: coprocess_object.proto
#pragma warning disable 1591, 0612, 3021
#region Designer generated code

using pb = global::Google.Protobuf;
using pbc = global::Google.Protobuf.Collections;
using pbr = global::Google.Protobuf.Reflection;
using scg = global::System.Collections.Generic;
namespace Coprocess {

  /// <summary>Holder for reflection information generated from coprocess_object.proto</summary>
  public static partial class CoprocessObjectReflection {

    #region Descriptor
    /// <summary>File descriptor for coprocess_object.proto</summary>
    public static pbr::FileDescriptor Descriptor {
      get { return descriptor; }
    }
    private static pbr::FileDescriptor descriptor;

    static CoprocessObjectReflection() {
      byte[] descriptorData = global::System.Convert.FromBase64String(
          string.Concat(
            "ChZjb3Byb2Nlc3Nfb2JqZWN0LnByb3RvEgljb3Byb2Nlc3MaI2NvcHJvY2Vz",
            "c19taW5pX3JlcXVlc3Rfb2JqZWN0LnByb3RvGh1jb3Byb2Nlc3Nfc2Vzc2lv",
            "bl9zdGF0ZS5wcm90bxoWY29wcm9jZXNzX2NvbW1vbi5wcm90byLYAgoGT2Jq",
            "ZWN0EiYKCWhvb2tfdHlwZRgBIAEoDjITLmNvcHJvY2Vzcy5Ib29rVHlwZRIR",
            "Cglob29rX25hbWUYAiABKAkSLQoHcmVxdWVzdBgDIAEoCzIcLmNvcHJvY2Vz",
            "cy5NaW5pUmVxdWVzdE9iamVjdBIoCgdzZXNzaW9uGAQgASgLMhcuY29wcm9j",
            "ZXNzLlNlc3Npb25TdGF0ZRIxCghtZXRhZGF0YRgFIAMoCzIfLmNvcHJvY2Vz",
            "cy5PYmplY3QuTWV0YWRhdGFFbnRyeRIpCgRzcGVjGAYgAygLMhsuY29wcm9j",
            "ZXNzLk9iamVjdC5TcGVjRW50cnkaLwoNTWV0YWRhdGFFbnRyeRILCgNrZXkY",
            "ASABKAkSDQoFdmFsdWUYAiABKAk6AjgBGisKCVNwZWNFbnRyeRILCgNrZXkY",
            "ASABKAkSDQoFdmFsdWUYAiABKAk6AjgBIhgKBUV2ZW50Eg8KB3BheWxvYWQY",
            "ASABKAkiDAoKRXZlbnRSZXBseTJ8CgpEaXNwYXRjaGVyEjIKCERpc3BhdGNo",
            "EhEuY29wcm9jZXNzLk9iamVjdBoRLmNvcHJvY2Vzcy5PYmplY3QiABI6Cg1E",
            "aXNwYXRjaEV2ZW50EhAuY29wcm9jZXNzLkV2ZW50GhUuY29wcm9jZXNzLkV2",
            "ZW50UmVwbHkiAGIGcHJvdG8z"));
      descriptor = pbr::FileDescriptor.FromGeneratedCode(descriptorData,
          new pbr::FileDescriptor[] { global::Coprocess.CoprocessMiniRequestObjectReflection.Descriptor, global::Coprocess.CoprocessSessionStateReflection.Descriptor, global::Coprocess.CoprocessCommonReflection.Descriptor, },
          new pbr::GeneratedClrTypeInfo(null, new pbr::GeneratedClrTypeInfo[] {
            new pbr::GeneratedClrTypeInfo(typeof(global::Coprocess.Object), global::Coprocess.Object.Parser, new[]{ "HookType", "HookName", "Request", "Session", "Metadata", "Spec" }, null, null, new pbr::GeneratedClrTypeInfo[] { null, null, }),
            new pbr::GeneratedClrTypeInfo(typeof(global::Coprocess.Event), global::Coprocess.Event.Parser, new[]{ "Payload" }, null, null, null),
            new pbr::GeneratedClrTypeInfo(typeof(global::Coprocess.EventReply), global::Coprocess.EventReply.Parser, null, null, null, null)
          }));
    }
    #endregion

  }
  #region Messages
  public sealed partial class Object : pb::IMessage<Object> {
    private static readonly pb::MessageParser<Object> _parser = new pb::MessageParser<Object>(() => new Object());
    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public static pb::MessageParser<Object> Parser { get { return _parser; } }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public static pbr::MessageDescriptor Descriptor {
      get { return global::Coprocess.CoprocessObjectReflection.Descriptor.MessageTypes[0]; }
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    pbr::MessageDescriptor pb::IMessage.Descriptor {
      get { return Descriptor; }
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public Object() {
      OnConstruction();
    }

    partial void OnConstruction();

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public Object(Object other) : this() {
      hookType_ = other.hookType_;
      hookName_ = other.hookName_;
      Request = other.request_ != null ? other.Request.Clone() : null;
      Session = other.session_ != null ? other.Session.Clone() : null;
      metadata_ = other.metadata_.Clone();
      spec_ = other.spec_.Clone();
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public Object Clone() {
      return new Object(this);
    }

    /// <summary>Field number for the "hook_type" field.</summary>
    public const int HookTypeFieldNumber = 1;
    private global::Coprocess.HookType hookType_ = 0;
    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public global::Coprocess.HookType HookType {
      get { return hookType_; }
      set {
        hookType_ = value;
      }
    }

    /// <summary>Field number for the "hook_name" field.</summary>
    public const int HookNameFieldNumber = 2;
    private string hookName_ = "";
    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public string HookName {
      get { return hookName_; }
      set {
        hookName_ = pb::ProtoPreconditions.CheckNotNull(value, "value");
      }
    }

    /// <summary>Field number for the "request" field.</summary>
    public const int RequestFieldNumber = 3;
    private global::Coprocess.MiniRequestObject request_;
    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public global::Coprocess.MiniRequestObject Request {
      get { return request_; }
      set {
        request_ = value;
      }
    }

    /// <summary>Field number for the "session" field.</summary>
    public const int SessionFieldNumber = 4;
    private global::Coprocess.SessionState session_;
    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public global::Coprocess.SessionState Session {
      get { return session_; }
      set {
        session_ = value;
      }
    }

    /// <summary>Field number for the "metadata" field.</summary>
    public const int MetadataFieldNumber = 5;
    private static readonly pbc::MapField<string, string>.Codec _map_metadata_codec
        = new pbc::MapField<string, string>.Codec(pb::FieldCodec.ForString(10), pb::FieldCodec.ForString(18), 42);
    private readonly pbc::MapField<string, string> metadata_ = new pbc::MapField<string, string>();
    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public pbc::MapField<string, string> Metadata {
      get { return metadata_; }
    }

    /// <summary>Field number for the "spec" field.</summary>
    public const int SpecFieldNumber = 6;
    private static readonly pbc::MapField<string, string>.Codec _map_spec_codec
        = new pbc::MapField<string, string>.Codec(pb::FieldCodec.ForString(10), pb::FieldCodec.ForString(18), 50);
    private readonly pbc::MapField<string, string> spec_ = new pbc::MapField<string, string>();
    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public pbc::MapField<string, string> Spec {
      get { return spec_; }
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public override bool Equals(object other) {
      return Equals(other as Object);
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public bool Equals(Object other) {
      if (ReferenceEquals(other, null)) {
        return false;
      }
      if (ReferenceEquals(other, this)) {
        return true;
      }
      if (HookType != other.HookType) return false;
      if (HookName != other.HookName) return false;
      if (!object.Equals(Request, other.Request)) return false;
      if (!object.Equals(Session, other.Session)) return false;
      if (!Metadata.Equals(other.Metadata)) return false;
      if (!Spec.Equals(other.Spec)) return false;
      return true;
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public override int GetHashCode() {
      int hash = 1;
      if (HookType != 0) hash ^= HookType.GetHashCode();
      if (HookName.Length != 0) hash ^= HookName.GetHashCode();
      if (request_ != null) hash ^= Request.GetHashCode();
      if (session_ != null) hash ^= Session.GetHashCode();
      hash ^= Metadata.GetHashCode();
      hash ^= Spec.GetHashCode();
      return hash;
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public override string ToString() {
      return pb::JsonFormatter.ToDiagnosticString(this);
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public void WriteTo(pb::CodedOutputStream output) {
      if (HookType != 0) {
        output.WriteRawTag(8);
        output.WriteEnum((int) HookType);
      }
      if (HookName.Length != 0) {
        output.WriteRawTag(18);
        output.WriteString(HookName);
      }
      if (request_ != null) {
        output.WriteRawTag(26);
        output.WriteMessage(Request);
      }
      if (session_ != null) {
        output.WriteRawTag(34);
        output.WriteMessage(Session);
      }
      metadata_.WriteTo(output, _map_metadata_codec);
      spec_.WriteTo(output, _map_spec_codec);
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public int CalculateSize() {
      int size = 0;
      if (HookType != 0) {
        size += 1 + pb::CodedOutputStream.ComputeEnumSize((int) HookType);
      }
      if (HookName.Length != 0) {
        size += 1 + pb::CodedOutputStream.ComputeStringSize(HookName);
      }
      if (request_ != null) {
        size += 1 + pb::CodedOutputStream.ComputeMessageSize(Request);
      }
      if (session_ != null) {
        size += 1 + pb::CodedOutputStream.ComputeMessageSize(Session);
      }
      size += metadata_.CalculateSize(_map_metadata_codec);
      size += spec_.CalculateSize(_map_spec_codec);
      return size;
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public void MergeFrom(Object other) {
      if (other == null) {
        return;
      }
      if (other.HookType != 0) {
        HookType = other.HookType;
      }
      if (other.HookName.Length != 0) {
        HookName = other.HookName;
      }
      if (other.request_ != null) {
        if (request_ == null) {
          request_ = new global::Coprocess.MiniRequestObject();
        }
        Request.MergeFrom(other.Request);
      }
      if (other.session_ != null) {
        if (session_ == null) {
          session_ = new global::Coprocess.SessionState();
        }
        Session.MergeFrom(other.Session);
      }
      metadata_.Add(other.metadata_);
      spec_.Add(other.spec_);
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public void MergeFrom(pb::CodedInputStream input) {
      uint tag;
      while ((tag = input.ReadTag()) != 0) {
        switch(tag) {
          default:
            input.SkipLastField();
            break;
          case 8: {
            hookType_ = (global::Coprocess.HookType) input.ReadEnum();
            break;
          }
          case 18: {
            HookName = input.ReadString();
            break;
          }
          case 26: {
            if (request_ == null) {
              request_ = new global::Coprocess.MiniRequestObject();
            }
            input.ReadMessage(request_);
            break;
          }
          case 34: {
            if (session_ == null) {
              session_ = new global::Coprocess.SessionState();
            }
            input.ReadMessage(session_);
            break;
          }
          case 42: {
            metadata_.AddEntriesFrom(input, _map_metadata_codec);
            break;
          }
          case 50: {
            spec_.AddEntriesFrom(input, _map_spec_codec);
            break;
          }
        }
      }
    }

  }

  public sealed partial class Event : pb::IMessage<Event> {
    private static readonly pb::MessageParser<Event> _parser = new pb::MessageParser<Event>(() => new Event());
    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public static pb::MessageParser<Event> Parser { get { return _parser; } }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public static pbr::MessageDescriptor Descriptor {
      get { return global::Coprocess.CoprocessObjectReflection.Descriptor.MessageTypes[1]; }
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    pbr::MessageDescriptor pb::IMessage.Descriptor {
      get { return Descriptor; }
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public Event() {
      OnConstruction();
    }

    partial void OnConstruction();

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public Event(Event other) : this() {
      payload_ = other.payload_;
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public Event Clone() {
      return new Event(this);
    }

    /// <summary>Field number for the "payload" field.</summary>
    public const int PayloadFieldNumber = 1;
    private string payload_ = "";
    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public string Payload {
      get { return payload_; }
      set {
        payload_ = pb::ProtoPreconditions.CheckNotNull(value, "value");
      }
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public override bool Equals(object other) {
      return Equals(other as Event);
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public bool Equals(Event other) {
      if (ReferenceEquals(other, null)) {
        return false;
      }
      if (ReferenceEquals(other, this)) {
        return true;
      }
      if (Payload != other.Payload) return false;
      return true;
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public override int GetHashCode() {
      int hash = 1;
      if (Payload.Length != 0) hash ^= Payload.GetHashCode();
      return hash;
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public override string ToString() {
      return pb::JsonFormatter.ToDiagnosticString(this);
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public void WriteTo(pb::CodedOutputStream output) {
      if (Payload.Length != 0) {
        output.WriteRawTag(10);
        output.WriteString(Payload);
      }
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public int CalculateSize() {
      int size = 0;
      if (Payload.Length != 0) {
        size += 1 + pb::CodedOutputStream.ComputeStringSize(Payload);
      }
      return size;
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public void MergeFrom(Event other) {
      if (other == null) {
        return;
      }
      if (other.Payload.Length != 0) {
        Payload = other.Payload;
      }
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public void MergeFrom(pb::CodedInputStream input) {
      uint tag;
      while ((tag = input.ReadTag()) != 0) {
        switch(tag) {
          default:
            input.SkipLastField();
            break;
          case 10: {
            Payload = input.ReadString();
            break;
          }
        }
      }
    }

  }

  public sealed partial class EventReply : pb::IMessage<EventReply> {
    private static readonly pb::MessageParser<EventReply> _parser = new pb::MessageParser<EventReply>(() => new EventReply());
    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public static pb::MessageParser<EventReply> Parser { get { return _parser; } }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public static pbr::MessageDescriptor Descriptor {
      get { return global::Coprocess.CoprocessObjectReflection.Descriptor.MessageTypes[2]; }
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    pbr::MessageDescriptor pb::IMessage.Descriptor {
      get { return Descriptor; }
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public EventReply() {
      OnConstruction();
    }

    partial void OnConstruction();

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public EventReply(EventReply other) : this() {
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public EventReply Clone() {
      return new EventReply(this);
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public override bool Equals(object other) {
      return Equals(other as EventReply);
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public bool Equals(EventReply other) {
      if (ReferenceEquals(other, null)) {
        return false;
      }
      if (ReferenceEquals(other, this)) {
        return true;
      }
      return true;
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public override int GetHashCode() {
      int hash = 1;
      return hash;
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public override string ToString() {
      return pb::JsonFormatter.ToDiagnosticString(this);
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public void WriteTo(pb::CodedOutputStream output) {
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public int CalculateSize() {
      int size = 0;
      return size;
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public void MergeFrom(EventReply other) {
      if (other == null) {
        return;
      }
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public void MergeFrom(pb::CodedInputStream input) {
      uint tag;
      while ((tag = input.ReadTag()) != 0) {
        switch(tag) {
          default:
            input.SkipLastField();
            break;
        }
      }
    }

  }

  #endregion

}

#endregion Designer generated code
