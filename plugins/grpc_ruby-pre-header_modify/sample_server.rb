require './tyk/dispatcher'
require './tyk/coprocess_session_state_pb'
require 'json'

class SampleServer < Coprocess::Dispatcher::Service
  # Implements a dynamic dispatcher for CP objects, this class should provide
  # methods for your hooks (see MyPreMiddleware).
  def dispatch(coprocess_object, _unused_call)
    begin
      if !coprocess_object.hook_name.nil?
        coprocess_object = self.send(coprocess_object.hook_name,
				     coprocess_object)
	dispatch_event(coprocess_object, _unused_call)
      else
        raise Coprocess::Dispatcher::HookNotImplemented
      end
    rescue Coprocess::Dispatcher::HookNotImplemented
      puts "Hook not implemented: #{coprocess_object.hook_name}"
    rescue Exception => e
      puts "Couldn't dispatch: #{e}"
    end
    return coprocess_object
  end

  # Implements an event dispatcher.
  def dispatch_event(event_wrapper, _unused_call)
    event = JSON.parse(event_wrapper.to_json)
    puts "Receiving object: #{event}"
    return Coprocess::EventReply.new
  end

  def MyPreMiddleware(coprocess_object)
    coprocess_object.request.set_headers["rubyheader"] = "rubyvalue"
    puts "Calling MyPreMiddleware"
    return coprocess_object
  end

  def MyPostMiddleware(coprocess_object)
    coprocess_object.request.set_headers["rubyheader"] = "rubyvalue"
    puts "Calling MyPostMiddleware"
    return coprocess_object
  end
end

def main
  s = GRPC::RpcServer.new
  s.add_http2_port('0.0.0.0:5555', :this_port_is_insecure)
  s.handle(SampleServer)
  s.run_till_terminated
end

main
