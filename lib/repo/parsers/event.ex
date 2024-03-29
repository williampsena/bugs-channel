defmodule BugsChannel.Repo.Parsers.Event do
  @moduledoc """
  This module includes database document parsers or mappers for the Service entity.
  """

  use BugsChannel.Repo.Parsers.Base

  alias BugsChannel.Utils.Maps
  alias BugsChannel.Utils.Ecto, as: EctoUtils
  alias BugsChannel.Repo.Schemas, as: RepoSchemas

  @doc ~S"""
  Parse a map to event.

  ## Examples

      iex> BugsChannel.Repo.Parsers.Event.parse(nil)
      nil

      iex> BugsChannel.Repo.Parsers.Event.parse(nil, %{})
      nil

      iex> BugsChannel.Repo.Parsers.Event.parse(%{}, nil)
      {:error, :invalid_schema}

      iex> BugsChannel.Repo.Parsers.Event.parse(%{ "id" => "00000000-0000-0000-0000-000000000000",  "service_id" => "1", "platform" => "python", "title" => "FooException: Bar messages", "body" => "Bar messages", "stack_trace" => [%{"type" => "FooException", "value" => "Bar messages"}], "kind" => "error", "level" => "error", "origin" => "sentry", "tags" => [] }, %BugsChannel.Repo.Schemas.Event{})
      %BugsChannel.Repo.Schemas.Event{
        id: "00000000-0000-0000-0000-000000000000",
        service_id: "1",
        platform: "python",
        environment: nil,
        release: nil,
        server_name: nil,
        title: "FooException: Bar messages",
        body: "Bar messages",
        stack_trace: [%{"type" => "FooException", "value" => "Bar messages"}],
        kind: "error",
        level: "error",
        origin: :sentry,
        tags: [],
        extra_args: nil
      }

  """
  def parse(doc), do: parse(doc, %RepoSchemas.Event{})

  def parse(nil, _schema), do: nil

  def parse(doc, %RepoSchemas.Event{} = schema)
      when is_map(doc) and is_struct(schema) do
    params = %{
      id: "#{doc["_id"] || doc["id"]}",
      service_id: "#{doc["service_id"]}",
      platform: doc["platform"],
      environment: doc["environment"],
      release: doc["release"],
      server_name: doc["server_name"],
      title: doc["title"],
      body: doc["body"],
      stack_trace: doc["stack_trace"],
      kind: doc["kind"],
      level: doc["level"],
      origin: String.to_existing_atom(doc["origin"]),
      tags: doc["tags"],
      extra_args: doc["extra_args"]
    }

    EctoUtils.parse_document(schema, params)
  end

  def parse(_doc, _schema), do: {:error, :invalid_schema}

  @doc ~S"""
  Parse a document list to schemas.

  ## Examples
      iex> doc = %{ "id" => "00000000-0000-0000-0000-000000000000",  "service_id" => "1", "platform" => "python", "title" => "FooException: Bar messages", "body" => "Bar messages", "stack_trace" => [%{"type" => "FooException", "value" => "Bar messages"}], "kind" => "error", "level" => "error", "origin" => "sentry", "tags" => [] }
      ...> BugsChannel.Repo.Parsers.Event.parse_list([doc], %BugsChannel.Repo.Schemas.Event{})
      [%BugsChannel.Repo.Schemas.Event{id: "00000000-0000-0000-0000-000000000000", service_id: "1", meta_id: nil, platform: "python", environment: nil, release: nil, server_name: nil, title: "FooException: Bar messages", body: "Bar messages", stack_trace: [%{"type" => "FooException", "value" => "Bar messages"}], kind: "error", level: "error", origin: :sentry, tags: [], extra_args: nil, inserted_at: nil, updated_at: nil}]
  """
  def parse_list(docs, schema) when is_list(docs) and is_struct(schema) do
    Enum.map(docs, fn doc -> __MODULE__.parse(doc, schema) end)
  end
end
