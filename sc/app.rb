require 'sinatra'
require 'socket'
require 'net/http'
$stdout.sync = true

service = ENV['SERVICE']
destinations = ENV['DESTINATIONS']


def get_headers env
  env.select { |k,v| k.start_with? 'HTTP_'}.
    transform_keys { |k| k.sub(/^HTTP_/, '').split('_').map(&:capitalize).join('-') }
end

def start_service service
  pod_ip = ENV['MY_POD_IP'] ||  '0.0.0.0'
  puts "starting service #{service}, bind to #{pod_ip}"

  set :bind, pod_ip
  set :port, 7000
  get "/#{service}" do
    content_type :text
    h = get_headers request.env
    puts h
    puts "=> receiving call from #{request.ip} REMOTE_ADDR: #{request.env["REMOTE_ADDR"]} X-Forwarded-For:#{h["X-Forwarded-For"]}"
    "<= #{service}(#{Socket.gethostname} #{my_first_private_ipv4}) say hi"
  end
end

def call_service service
  puts "calling service #{service}:"
  uri = URI("http://#{service}:7000/#{service}")
  req = Net::HTTP::Get.new(uri)
  res = Net::HTTP.start(uri.hostname, uri.port) {|http|
    http.request(req)
  }
  puts res.body
end

def call_destinations destinations
  services = destinations.split ","
  loop do 
    services.each do |service|
      begin
      call_service service
      rescue => e
        puts "call service #{service} error: #{e}"
      end
    end
    sleep 1
  end
end

if service
  start_service service
elsif destinations
  call_destinations destinations
else
  puts "miss env service or destinations, exiting"
  exit 1
end

def my_first_private_ipv4
  ip = Socket.ip_address_list.detect{|intf| intf.ipv4_private?}
  ip && ip.ip_address
end

# def my_first_public_ipv4
#   Socket.ip_address_list.detect{|intf| intf.ipv4? and !intf.ipv4_loopback? and !intf.ipv4_multicast? and !intf.ipv4_private?}
# end
