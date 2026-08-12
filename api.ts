/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export enum EntityTaskStatus {
  TaskStatusTodo = "todo",
  TaskStatusInProgress = "in_progress",
  TaskStatusDone = "done",
}

export interface EntityTask {
  /** @example "2026-01-01T00:00:00Z" */
  created_at?: string;
  /** @example "Task description" */
  description?: string;
  /** @example "550e8400-e29b-41d4-a716-446655440000" */
  id?: string;
  /** @example "todo" */
  status?: EntityTaskStatus;
  /** @example "My task" */
  title?: string;
  /** @example "2026-01-01T00:00:00Z" */
  updated_at?: string;
  /** @example "550e8400-e29b-41d4-a716-446655440000" */
  user_id?: string;
}

export interface EntityTranslation {
  /** @example "en" */
  destination?: string;
  /** @example "текст для перевода" */
  original?: string;
  /** @example "auto" */
  source?: string;
  /** @example "text for translation" */
  translation?: string;
}

export interface EntityTranslationHistory {
  history?: EntityTranslation[];
}

export interface EntityUser {
  /** @example "2026-01-01T00:00:00Z" */
  created_at?: string;
  /** @example "john@example.com" */
  email?: string;
  /** @example "550e8400-e29b-41d4-a716-446655440000" */
  id?: string;
  /** @example "2026-01-01T00:00:00Z" */
  updated_at?: string;
  /** @example "johndoe" */
  username?: string;
}

export interface V1CreateTask {
  /**
   * @maxLength 1000
   * @example "Task description"
   */
  description?: string;
  /**
   * @maxLength 255
   * @example "My task"
   */
  title: string;
}

export interface V1Error {
  /** @example "message" */
  error?: string;
}

export interface V1Login {
  /** @example "john@example.com" */
  email: string;
  /** @example "secret123" */
  password: string;
}

export interface V1Register {
  /** @example "john@example.com" */
  email: string;
  /**
   * @minLength 6
   * @example "secret123"
   */
  password: string;
  /**
   * @minLength 3
   * @maxLength 255
   * @example "johndoe"
   */
  username: string;
}

export interface V1TaskList {
  tasks?: EntityTask[];
  /** @example 42 */
  total?: number;
}

export interface V1Token {
  /** @example "eyJhbGciOiJIUzI1NiIs..." */
  token?: string;
}

export interface V1TransitionTask {
  /** @example "in_progress" */
  status: "todo" | "in_progress" | "done";
}

export interface V1Translate {
  /** @example "en" */
  destination: string;
  /** @example "текст для перевода" */
  original: string;
  /** @example "auto" */
  source: string;
}

export interface V1UpdateTask {
  /**
   * @maxLength 1000
   * @example "Updated description"
   */
  description?: string;
  /**
   * @maxLength 255
   * @example "Updated task"
   */
  title: string;
}

export type QueryParamsType = Record<string | number, any>;
export type ResponseFormat = keyof Omit<Body, "body" | "bodyUsed">;

export interface FullRequestParams extends Omit<RequestInit, "body"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseFormat;
  /** request body */
  body?: unknown;
  /** base url */
  baseUrl?: string;
  /** request cancellation token */
  cancelToken?: CancelToken;
}

export type RequestParams = Omit<
  FullRequestParams,
  "body" | "method" | "query" | "path"
>;

export interface ApiConfig<SecurityDataType = unknown> {
  baseUrl?: string;
  baseApiParams?: Omit<RequestParams, "baseUrl" | "cancelToken" | "signal">;
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<RequestParams | void> | RequestParams | void;
  customFetch?: typeof fetch;
}

export interface HttpResponse<D extends unknown, E extends unknown = unknown>
  extends Response {
  data: D;
  error: E;
}

type CancelToken = Symbol | string | number;

export enum ContentType {
  Json = "application/json",
  JsonApi = "application/vnd.api+json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
  Text = "text/plain",
}

export class HttpClient<SecurityDataType = unknown> {
  public baseUrl: string = "//localhost:8080/v1";
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private abortControllers = new Map<CancelToken, AbortController>();
  private customFetch = (...fetchParams: Parameters<typeof fetch>) =>
    fetch(...fetchParams);

  private baseApiParams: RequestParams = {
    credentials: "same-origin",
    headers: {},
    redirect: "follow",
    referrerPolicy: "no-referrer",
  };

  constructor(apiConfig: ApiConfig<SecurityDataType> = {}) {
    Object.assign(this, apiConfig);
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected encodeQueryParam(key: string, value: any) {
    const encodedKey = encodeURIComponent(key);
    return `${encodedKey}=${encodeURIComponent(typeof value === "number" ? value : `${value}`)}`;
  }

  protected addQueryParam(query: QueryParamsType, key: string) {
    return this.encodeQueryParam(key, query[key]);
  }

  protected addArrayQueryParam(query: QueryParamsType, key: string) {
    const value = query[key];
    return value.map((v: any) => this.encodeQueryParam(key, v)).join("&");
  }

  protected toQueryString(rawQuery?: QueryParamsType): string {
    const query = rawQuery || {};
    const keys = Object.keys(query).filter(
      (key) => "undefined" !== typeof query[key],
    );
    return keys
      .map((key) =>
        Array.isArray(query[key])
          ? this.addArrayQueryParam(query, key)
          : this.addQueryParam(query, key),
      )
      .join("&");
  }

  protected addQueryParams(rawQuery?: QueryParamsType): string {
    const queryString = this.toQueryString(rawQuery);
    return queryString ? `?${queryString}` : "";
  }

  private contentFormatters: Record<ContentType, (input: any) => any> = {
    [ContentType.Json]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string")
        ? JSON.stringify(input)
        : input,
    [ContentType.JsonApi]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string")
        ? JSON.stringify(input)
        : input,
    [ContentType.Text]: (input: any) =>
      input !== null && typeof input !== "string"
        ? JSON.stringify(input)
        : input,
    [ContentType.FormData]: (input: any) => {
      if (input instanceof FormData) {
        return input;
      }

      return Object.keys(input || {}).reduce((formData, key) => {
        const property = input[key];
        formData.append(
          key,
          property instanceof Blob
            ? property
            : typeof property === "object" && property !== null
              ? JSON.stringify(property)
              : `${property}`,
        );
        return formData;
      }, new FormData());
    },
    [ContentType.UrlEncoded]: (input: any) => this.toQueryString(input),
  };

  protected mergeRequestParams(
    params1: RequestParams,
    params2?: RequestParams,
  ): RequestParams {
    return {
      ...this.baseApiParams,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.baseApiParams.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected createAbortSignal = (
    cancelToken: CancelToken,
  ): AbortSignal | undefined => {
    if (this.abortControllers.has(cancelToken)) {
      const abortController = this.abortControllers.get(cancelToken);
      if (abortController) {
        return abortController.signal;
      }
      return void 0;
    }

    const abortController = new AbortController();
    this.abortControllers.set(cancelToken, abortController);
    return abortController.signal;
  };

  public abortRequest = (cancelToken: CancelToken) => {
    const abortController = this.abortControllers.get(cancelToken);

    if (abortController) {
      abortController.abort();
      this.abortControllers.delete(cancelToken);
    }
  };

  public request = async <T = any, E = any>({
    body,
    secure,
    path,
    type,
    query,
    format,
    baseUrl,
    cancelToken,
    ...params
  }: FullRequestParams): Promise<HttpResponse<T, E>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.baseApiParams.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const queryString = query && this.toQueryString(query);
    const payloadFormatter = this.contentFormatters[type || ContentType.Json];
    const responseFormat = format || requestParams.format;

    return this.customFetch(
      `${baseUrl || this.baseUrl || ""}${path}${queryString ? `?${queryString}` : ""}`,
      {
        ...requestParams,
        headers: {
          ...(requestParams.headers || {}),
          ...(type && type !== ContentType.FormData
            ? { "Content-Type": type }
            : {}),
        },
        signal:
          (cancelToken
            ? this.createAbortSignal(cancelToken)
            : requestParams.signal) || null,
        body:
          typeof body === "undefined" || body === null
            ? null
            : payloadFormatter(body),
      },
    ).then(async (response) => {
      const r = response as HttpResponse<T, E>;
      r.data = null as unknown as T;
      r.error = null as unknown as E;

      const responseToParse = responseFormat ? response.clone() : response;
      const data = !responseFormat
        ? r
        : await responseToParse[responseFormat]()
            .then((data) => {
              if (r.ok) {
                r.data = data;
              } else {
                r.error = data;
              }
              return r;
            })
            .catch((e) => {
              r.error = e;
              return r;
            });

      if (cancelToken) {
        this.abortControllers.delete(cancelToken);
      }

      if (!response.ok) throw data;
      return data;
    });
  };
}

/**
 * @title Go Clean Template API
 * @version 1.0
 * @baseUrl //localhost:8080/v1
 * @contact
 *
 * Multi-domain clean architecture template with translation, user, and task management
 */
export class Api<
  SecurityDataType extends unknown,
> extends HttpClient<SecurityDataType> {
  auth = {
    /**
     * @description Authenticate user and get JWT token
     *
     * @tags auth
     * @name Login
     * @summary Login
     * @request POST:/auth/login
     */
    login: (request: V1Login, params: RequestParams = {}) =>
      this.request<V1Token, V1Error>({
        path: `/auth/login`,
        method: "POST",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Register a new user
     *
     * @tags auth
     * @name Register
     * @summary Register
     * @request POST:/auth/register
     */
    register: (request: V1Register, params: RequestParams = {}) =>
      this.request<EntityUser, V1Error>({
        path: `/auth/register`,
        method: "POST",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  tasks = {
    /**
     * @description List tasks for the current user with optional filtering
     *
     * @tags tasks
     * @name ListTasks
     * @summary List tasks
     * @request GET:/tasks
     * @secure
     */
    listTasks: (
      query?: {
        /** Filter by status */
        status?: "todo" | "in_progress" | "done";
        /**
         * Limit
         * @default 10
         */
        limit?: number;
        /**
         * Offset
         * @default 0
         */
        offset?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<V1TaskList, V1Error>({
        path: `/tasks`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Create a new task for the current user
     *
     * @tags tasks
     * @name CreateTask
     * @summary Create task
     * @request POST:/tasks
     * @secure
     */
    createTask: (request: V1CreateTask, params: RequestParams = {}) =>
      this.request<EntityTask, V1Error>({
        path: `/tasks`,
        method: "POST",
        body: request,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Delete a task by ID
     *
     * @tags tasks
     * @name DeleteTask
     * @summary Delete task
     * @request DELETE:/tasks/{id}
     * @secure
     */
    deleteTask: (id: string, params: RequestParams = {}) =>
      this.request<void, V1Error>({
        path: `/tasks/${id}`,
        method: "DELETE",
        secure: true,
        ...params,
      }),

    /**
     * @description Get a task by ID
     *
     * @tags tasks
     * @name GetTask
     * @summary Get task
     * @request GET:/tasks/{id}
     * @secure
     */
    getTask: (id: string, params: RequestParams = {}) =>
      this.request<EntityTask, V1Error>({
        path: `/tasks/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Update task title and description
     *
     * @tags tasks
     * @name UpdateTask
     * @summary Update task
     * @request PUT:/tasks/{id}
     * @secure
     */
    updateTask: (
      id: string,
      request: V1UpdateTask,
      params: RequestParams = {},
    ) =>
      this.request<EntityTask, V1Error>({
        path: `/tasks/${id}`,
        method: "PUT",
        body: request,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Change task status (todo -> in_progress -> done, or in_progress -> todo)
     *
     * @tags tasks
     * @name TransitionTask
     * @summary Transition task status
     * @request PATCH:/tasks/{id}/status
     * @secure
     */
    transitionTask: (
      id: string,
      request: V1TransitionTask,
      params: RequestParams = {},
    ) =>
      this.request<EntityTask, V1Error>({
        path: `/tasks/${id}/status`,
        method: "PATCH",
        body: request,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  translation = {
    /**
     * @description Translate a text
     *
     * @tags translation
     * @name DoTranslate
     * @summary Translate
     * @request POST:/translation/do-translate
     * @secure
     */
    doTranslate: (request: V1Translate, params: RequestParams = {}) =>
      this.request<EntityTranslation, V1Error>({
        path: `/translation/do-translate`,
        method: "POST",
        body: request,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Show all translation history for current user
     *
     * @tags translation
     * @name History
     * @summary Show history
     * @request GET:/translation/history
     * @secure
     */
    history: (params: RequestParams = {}) =>
      this.request<EntityTranslationHistory, V1Error>({
        path: `/translation/history`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),
  };
  user = {
    /**
     * @description Get current user profile
     *
     * @tags user
     * @name Profile
     * @summary Get profile
     * @request GET:/user/profile
     * @secure
     */
    profile: (params: RequestParams = {}) =>
      this.request<EntityUser, V1Error>({
        path: `/user/profile`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),
  };
}
