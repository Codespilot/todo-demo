"use client";

import { useEffect, useState } from "react";
import { client } from "../lib/client";
import { Task } from "../gen/todo_pb";

export default function Home() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [text, setText] = useState("");

  const loadTasks = async () => {
    const res = await client.getTasks({});
    setTasks(res.tasks);
  };

  const addTask = async () => {
    if (!text.trim()) return;
    await client.addTask({ text });
    setText("");
    await loadTasks();
  };

  useEffect(() => {
    loadTasks();
  }, []);

  return (
    <main style={{ padding: "2rem" }}>
      <h1>待办事项</h1>
      <div>
        <input
          value={text}
          onChange={(e) => setText(e.target.value)}
          placeholder="输入任务"
        />
        <button onClick={addTask}>添加</button>
      </div>
      <ul>
        {tasks.map((t) => (
          <li key={t.id}>{t.text}</li>
        ))}
      </ul>
    </main>
  );
}