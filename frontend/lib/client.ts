import { createConnectTransport } from "@connectrpc/connect-web";
import { createClient } from "@connectrpc/connect";
import { TodoService } from "../gen/todo_pb";

const transport = createConnectTransport({
  baseUrl: "http://localhost:8000", // 后端地址
});

export const client = createClient(TodoService, transport);