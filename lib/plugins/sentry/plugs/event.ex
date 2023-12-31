defmodule BugsChannel.Plugins.Sentry.Plugs.Event do
  @moduledoc """
  This plug is in responsible for handling sentry issues
  """

  require Logger

  import BugsChannel.Plugs.Api
  import Plug.Conn

  alias BugsChannel.Events.RawEvent

  def init(options) do
    options
  end

  def call(conn, _opts) do
    action(conn, conn.method)
  end

  defp action(%Plug.Conn{params: %{"event_id" => event_id}} = conn, "POST") do
    if is_nil(event_id) do
      send_resp(conn, 204, "")
    else
      case RawEvent.publish(event_id, "sentry", conn.params) do
        :ok ->
          send_json_resp(conn, %{"event_id" => event_id})

        error ->
          Logger.error("❌ An error occurred while attempting to send raw events: #{inspect(error)}")

          send_unknown_error_resp(conn)
      end
    end
  end

  defp action(conn, "POST") do
    send_unprocessable_entity_resp(conn)
  end

  defp action(conn, _) do
    send_not_found_resp(conn)
  end
end
