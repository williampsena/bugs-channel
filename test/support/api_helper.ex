defmodule BugsChannel.Test.Support.ApiHelper do
  @moduledoc """
  This module includes various supports for asserting and dealing with Plug.Conn.
  """

  use ExUnit.Case, only: [assert: 1]

  def assert_conn(%Plug.Conn{} = conn, status_code, body) when is_integer(status_code) do
    assert conn.state == :sent
    assert conn.status == status_code

    if is_map(body) do
      assert Jason.decode!(conn.resp_body) == body
    else
      assert conn.resp_body == body
    end
  end

  def put_params(conn, params) do
    %{conn | params: params}
  end
end
