import { httpClient } from "./http";

export interface APIKey {
  id: string;
  name: string;
  scopes: string;
  created_at: string;
  key?: string; // Only present on creation
}

export interface CreateAPIKeyResponse {
  id: string;
  name: string;
  scopes: string;
  created_at: string;
  key: string; // The plain text key (only shown once)
}

export interface ListAPIKeysResponse {
  keys: APIKey[];
}

class APIKeyService {
  /**
   * Create a new API key
   */
  async createAPIKey(name: string): Promise<CreateAPIKeyResponse> {
    const response = await httpClient.post<CreateAPIKeyResponse>(
      "/api/api-keys",
      { name }
    );
    return response;
  }

  /**
   * Get all API keys for the current user
   */
  async listAPIKeys(): Promise<APIKey[]> {
    const response = await httpClient.get<ListAPIKeysResponse>("/api/api-keys");
    return response.keys;
  }

  /**
   * Revoke (delete) an API key
   */
  async revokeAPIKey(keyId: string): Promise<void> {
    await httpClient.delete(`/api/api-keys/${keyId}`);
  }

  /**
   * Mask an API key for display (show only first and last 4 characters)
   */
  maskAPIKey(key: string): string {
    if (key.length <= 8) {
      return "*".repeat(key.length);
    }
    return key.slice(0, 4) + "*".repeat(key.length - 8) + key.slice(-4);
  }
}

export const apiKeyService = new APIKeyService();
