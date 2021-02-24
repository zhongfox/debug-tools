require 'sinatra'
require 'socket'
require 'net/http'
$stdout.sync = true

def get_headers env
  env.select { |k,v| k.start_with? 'HTTP_'}.
    transform_keys { |k| k.sub(/^HTTP_/, '').split('_').map(&:capitalize).join('-') }
end

def start_http_service(service, port)
  puts "starting http service #{service} on port #{port}"

  set :bind, pod_ip
  set :port, port
  get "/#{service}" do
    content_type :text
    h = get_headers request.env
    puts h
    puts "=> Receiving call from #{request.ip} REMOTE_ADDR: #{request.env["REMOTE_ADDR"]} X-Forwarded-For:#{h["X-Forwarded-For"]}"
    "<= #{service}(#{Socket.gethostname} #{pod_ip}) say hi"
  end
end

def start_tcp_service(service, port)
  puts "starting tcp service #{service} on port #{port}"
  server = TCPServer.new(port)

  loop do
    Thread.start(server.accept) do |connection|
      sock_domain, remote_port, remote_hostname, remote_ip = connection.peeraddr
      client = "#{remote_ip}:#{remote_port}"
      puts "=>#{Time.now} Connection from #{client}\n"

      while line = connection.gets
        break if line =~ /quit/
        puts "=>Received from #{client}: #{line}"
        connection.puts "Server(#{pod_ip}) received\n"
      end
      puts "Closing the connection #{client}\n"
      connection.close
    end
  end
end

def call_http_service(service, port)
  puts "calling service #{service}:"
  uri = URI("http://#{service}:#{port}/#{service}")
  req = Net::HTTP::Get.new(uri)
  res = Net::HTTP.start(uri.hostname, uri.port) {|http|
    http.request(req)
  }
  puts res.body
end

def call_http_services(destinations, port)
  services = destinations.split ","
  loop do 
    services.each do |service|
      begin
      call_http_service(service, port)
      rescue => e
        puts "call service #{service} error: #{e}"
      end
    end
    sleep 1
  end
end

def pod_ip
   if ENV['POD_IP']
     return ENV['POD_IP']
   end

  ip = Socket.ip_address_list.detect{|intf| intf.ipv4_private?}
  ip && ip.ip_address
end

http_port = ENV['HTTP_PORT'] || 7000
tcp_port = ENV['TCP_PORT'] || 4000
service = ENV['SERVICE']
destinations = ENV['DESTINATIONS']

if service
  Thread.new { start_tcp_service(service, tcp_port) }
  start_http_service(service, http_port)
elsif destinations
  call_http_services(destinations, http_port)
else
  puts "miss env service or destinations, exiting"
  exit 1
end
