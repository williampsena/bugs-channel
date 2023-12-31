defmodule BugsChannel.Repo.DBless.Service do
  @moduledoc """
  This module contains the DBLess mode service repository.
  """

  alias BugsChannel.Settings.Context
  alias BugsChannel.Repo.Schemas, as: RepoSchemas

  @spec get(String.t()) :: RepoSchemas.Service
  def get(id) do
    do_query(&(&1.id == id))
  end

  @spec get_by_auth_key(String.t()) :: RepoSchemas.Service
  def get_by_auth_key(auth_key) do
    do_query(fn service ->
      Enum.any?(service.auth_keys, fn key ->
        key.key == auth_key &&
          key.disabled == false &&
          (is_nil(key.expired_at) ||
             :os.system_time(:seconds) < key.expired_at)
      end)
    end)
  end

  defp do_query(query) do
    config_file = Context.manager().get_config_file()
    Enum.find(config_file.services, query)
  end
end
