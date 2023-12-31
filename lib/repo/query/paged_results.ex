defmodule BugsChannel.Repo.Query.PagedResults do
  @moduledoc """
  The module represents paged results
  """
  alias BugsChannel.Repo.Query.QueryCursor

  defstruct ~w(data meta local)a

  @doc ~S"""
  Build list with paged results metadata

  ## Examples

      iex> BugsChannel.Repo.Query.PagedResults.build(~w(foo bar), %BugsChannel.Repo.Query.QueryCursor{page: 0, offset: 0, limit: 10})
      %BugsChannel.Repo.Query.PagedResults{data: ["foo", "bar"],  local: %{empty: false, next_page:  %BugsChannel.Repo.Query.QueryCursor{offset: 10, limit: 10, page: 1}}, meta: %{ page: 0, offset: 0, limit: 10, count: 2 }}

      iex> BugsChannel.Repo.Query.PagedResults.build([], %BugsChannel.Repo.Query.QueryCursor{page: 0, offset: 0, limit: 10})
      %BugsChannel.Repo.Query.PagedResults{data: [], local: %{empty: true}, meta: %{ page: 0, offset: 0, limit: 10, count: 0 }}
  """
  def build(results, %QueryCursor{} = query_cursor) when is_list(results) do
    count = Kernel.length(results)

    paged_results = %__MODULE__{
      data: results,
      meta: %{
        page: query_cursor.page,
        offset: query_cursor.offset,
        limit: query_cursor.limit,
        count: count
      },
      local: %{
        empty: count == 0
      }
    }

    if count > 0 do
      put_in(
        paged_results.local[:next_page],
        query_cursor |> QueryCursor.build_next()
      )
    else
      paged_results
    end
  end

  @doc ~S"""
  Build list with paged results metadata and next url

  ## Examples

      iex> BugsChannel.Repo.Query.PagedResults.build_with_next(~w(foo bar), "/foo", %{"bar" => "foo"}, %BugsChannel.Repo.Query.QueryCursor{page: 0, offset: 0, limit: 10})
      %BugsChannel.Repo.Query.PagedResults{
        data: ["foo", "bar"],
        local: %{empty: false, next_page: %BugsChannel.Repo.Query.QueryCursor{offset: 10, limit: 10, page: 1}},
        meta: %{limit: 10, offset: 0, page: 0, count: 2, next_url: "/foo?index=1&limit=10&bar=foo"}
      }

      iex> BugsChannel.Repo.Query.PagedResults.build_with_next(%BugsChannel.Repo.Query.PagedResults{ data: ["foo", "bar"], meta: %{limit: 10, offset: 0, page: 0 } }, "/foo", %{"bar" => "foo"}, %BugsChannel.Repo.Query.QueryCursor{page: 0, offset: 0, limit: 10})
      %BugsChannel.Repo.Query.PagedResults{
        data: ["foo", "bar"],
        meta: %{limit: 10, offset: 0, page: 0, next_url: "/foo?index=1&limit=10&bar=foo"}
      }

      iex> BugsChannel.Repo.Query.PagedResults.build_with_next(%BugsChannel.Repo.Query.PagedResults{ data: ["foo", "bar"], meta: %{limit: 10, offset: 0, page: 0 } }, "/foo", %{"bar" => "foo"})
      %BugsChannel.Repo.Query.PagedResults{
        data: ["foo", "bar"],
        meta: %{limit: 10, offset: 0, page: 0, next_url: "/foo?index=1&limit=10&bar=foo"}
      }

      iex> BugsChannel.Repo.Query.PagedResults.build_with_next(%BugsChannel.Repo.Query.PagedResults{ data: ["foo", "bar"], local: %{empty: true}, meta: %{limit: 10, offset: 0, page: 0 } }, "/foo", %{"bar" => "foo"})
      %BugsChannel.Repo.Query.PagedResults{
        data: ["foo", "bar"],
        local: %{empty: true},
        meta: %{limit: 10, offset: 0, page: 0}
      }
  """

  def build_with_next(
        %__MODULE__{} = paged_results,
        uri_path,
        query_params,
        %QueryCursor{} = query_cursor
      )
      when is_binary(uri_path) and is_map(query_params) do
    case paged_results do
      %__MODULE__{local: %{empty: true}} ->
        paged_results

      _ ->
        put_in(
          paged_results.meta[:next_url],
          QueryCursor.build_next_url(query_cursor, uri_path, query_params)
        )
    end
  end

  def build_with_next(results, uri_path, query_params, %QueryCursor{} = query_cursor)
      when is_list(results) and is_binary(uri_path) and is_map(query_params) do
    paged_results = build(results, query_cursor)

    build_with_next(paged_results, uri_path, query_params, query_cursor)
  end

  def build_with_next(
        %__MODULE__{meta: meta} = paged_results,
        uri_path,
        query_params
      )
      when is_binary(uri_path) and is_map(query_params) do
    build_with_next(
      paged_results,
      uri_path,
      query_params,
      QueryCursor.build(meta.page, meta.limit)
    )
  end
end
