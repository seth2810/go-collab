import React, { useState, useEffect, useCallback } from 'react';
import { Delta, Sources, DeltaStatic } from "quill";
import ReactQuill from 'react-quill';

import 'react-quill/dist/quill.snow.css';
import './App.css';

enum MessageTypes {
  SaveDocument = "save-document",
  LoadDocument = "load-document",
  UpdateDocument = "update-document",
  UpdateUserCount = "update-user-count",
}

interface Editor {
  getContents(index?: number, length?: number): DeltaStatic;
}

type SocketMessage = {
  type: MessageTypes;
  payload: any;
};

const modules = {
  toolbar: false,
}

const App: React.FunctionComponent = () => {
  const [readOnly, setReadOnly] = useState(true);
  const [value, setValue] = useState<Delta | string>("");
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [usersCount, setUsersCount] = useState<number | null>(null);

  const handleChange = useCallback((_: string, delta: Delta, sources: Sources, editor: Editor) => {
    if (sources !== "user" || !socket) {
      return;
    }

    const content = editor.getContents()

    const updateDocumentData = JSON.stringify({
      type: MessageTypes.UpdateDocument,
      payload: content,
    });

    socket.send(updateDocumentData);

    const saveDocumentData = JSON.stringify({
      type: MessageTypes.SaveDocument,
      payload: content,
    });

    socket.send(saveDocumentData);
  }, [socket])

  const handleMessage = useCallback((event: MessageEvent<string>) => {
    const message: SocketMessage = JSON.parse(event.data);

    switch (message.type) {
      case MessageTypes.LoadDocument:
        setValue(message.payload as Delta);
        setReadOnly(false);
        break;
      case MessageTypes.UpdateDocument:
        setValue(message.payload as Delta);
        break;
      case MessageTypes.UpdateUserCount:
        setUsersCount(message.payload as number);
        break;
    }
  }, [setReadOnly, setUsersCount]);

  useEffect(() => {
    const socket = new WebSocket(process.env.REACT_APP_BACKEND_URL ?? location.href);

    socket.addEventListener('message', handleMessage);

    setSocket(socket);

    return () => {
      socket.removeEventListener('message', handleMessage);
      socket.close();
    }
  }, [handleMessage]);

  return (
    <div className='container'>
      <h3>Users count: {usersCount}</h3>
      <ReactQuill
        theme='snow'
        value={value}
        modules={modules}
        readOnly={readOnly}
        onChange={handleChange}
        defaultValue="Loading..."
      />
    </div>
  );
};

export default App;
