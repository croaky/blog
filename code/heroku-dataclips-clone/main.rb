#!/usr/bin/env ruby

require "bundler/inline"
require "csv"

gemfile do
  source "https://rubygems.org"
  gem "haml"
  gem "pg"
  gem "puma"
  gem "sinatra"
end

enable :inline_templates

db = PG.connect(ENV.fetch("DATABASE_URL"))

get "/" do
  haml :index
end

post "/run" do
  @rows = db.exec(params["sql"] || "SELECT 1")
  haml :index
end

post "/export" do
  headers "Content-Type" => "text/csv"
  headers "Content-Disposition" => "attachment; filename=dataclip.csv"

  stream = StringIO.new

  db.copy_data("COPY (#{params["sql"] || "SELECT 1"}) TO STDOUT DELIMITER ',' CSV HEADER FORCE QUOTE *") do
    while (row = db.get_copy_data)
      stream.write(row)
    end
  end

  stream.string
end

Sinatra::Application.run!

__END__

@@ index
%form(method="POST")
  %textarea{name: "sql", rows: 10, cols: 80}= params["sql"]
  %br
  %button{type: "submit", formaction: "/run"} Run
  %button{type: "submit", formaction: "/export"} Export

- if @rows
  %table
    %thead
      %tr
        - @rows.fields.each do |field|
          %th= field
    %tbody
      - @rows.each do |row|
        %tr
          - row.each do |key, value|
            %td= value
