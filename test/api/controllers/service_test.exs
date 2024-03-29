defmodule BugsChannel.Api.Controllers.ServiceTest do
  use ExUnit.Case
  use Plug.Test

  import Mock
  import BugsChannel.Factories.Service
  import BugsChannel.Test.Support.ApiHelper

  alias BugsChannel.Utils.Maps
  alias BugsChannel.Repo
  alias BugsChannel.Api.Controllers.Service, as: ServiceController
  alias BugsChannel.Plugs.ErrorFallback

  setup do
    service = build(:service)

    {:ok, inserted_service} =
      Repo.Service.insert(service)

    [service: inserted_service, service_id: "#{inserted_service.id}"]
  end

  describe "GET /service/:id" do
    test "returns not found" do
      service_id = "657529b00000000000000000"

      conn =
        :get
        |> conn("/service/#{service_id}", "")
        |> ServiceController.show(%{"id" => service_id})

      assert_conn(conn, 404, "Oops! 👀")
    end

    test "returns a service", %{service: service, service_id: service_id} do
      conn =
        :get
        |> conn("/service/#{service_id}", "")
        |> ServiceController.show(%{"id" => service_id})

      service_map = service |> Maps.map_from_struct() |> Jason.encode!()

      assert_conn(conn, 200, service_map)
    end
  end

  describe "POST /service" do
    setup do
      service_params = %{
        "name" => "foo bar service",
        "platform" => "python",
        "teams" => ["1"],
        "settings" => %{"rate_limit" => 1},
        "auth_keys" => [%{"key" => "123"}]
      }

      [service_params: service_params]
    end

    test "with success", %{service_params: service_params} do
      conn =
        :post
        |> conn("/service", "")
        |> ServiceController.create(service_params)

      assert_conn(conn, 201, :skip)
      assert match?(%{"id" => _service_id}, Jason.decode!(conn.resp_body))
    end

    test "with validation error", %{service_id: service_id} do
      params = %{"id" => service_id}

      conn = conn(:post, "/service", "")

      conn =
        conn
        |> ServiceController.create(params)
        |> ErrorFallback.fallback(conn)

      assert_conn(
        conn,
        422,
        Jason.encode!(%{"error" => ["platform: can't be blank", "name: can't be blank"]})
      )
    end

    test "with error", %{service_params: params} do
      error = "invalid connection"

      with_mock(Repo.Service, [:passthrough], insert: fn _service -> {:error, error} end) do
        conn = conn(:post, "/service", "")

        conn =
          conn
          |> ServiceController.create(params)
          |> ErrorFallback.fallback(conn)

        assert_conn(conn, 500, %{"error" => error})
      end
    end
  end

  describe "PATCH /service" do
    setup context do
      service_params = %{
        "id" => context.service_id,
        "name" => "foo bar service #update"
      }

      [service_params: service_params]
    end

    test "with success", %{service_id: service_id, service_params: params} do
      conn =
        :patch
        |> conn("/service", %{"id" => service_id})
        |> ServiceController.update(params)

      assert_conn(conn, 204, "")
    end

    test "with validation error", %{service_id: service_id} do
      params = %{"id" => service_id, "foo" => "bar"}

      conn = conn(:patch, "/service", %{"id" => service_id})

      conn =
        conn
        |> ServiceController.update(params)
        |> ErrorFallback.fallback(conn)

      assert_conn(
        conn,
        422,
        Jason.encode!(%{"error" => "There are no fields to be updated."})
      )
    end

    test "with error", %{service_id: service_id, service_params: params} do
      error = "invalid connection"

      with_mock(Repo.Service, [:passthrough], update: fn _id, _service -> {:error, error} end) do
        conn = conn(:patch, "/service", %{"id" => service_id})

        conn =
          conn
          |> ServiceController.update(params)
          |> ErrorFallback.fallback(conn)

        assert_conn(conn, 500, %{"error" => error})
      end
    end
  end
end
