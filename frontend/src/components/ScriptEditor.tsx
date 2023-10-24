import React, { useContext, useEffect } from 'react';
import { CodeEditor } from '../ui/Editor';

import { Context } from '../state';
import { UpdateScriptFileContent } from '../../wailsjs/go/api/Api';

export const ScriptEditor: React.FC = () => {
  const ctx = useContext(Context);
  let fileId = ctx.state.activeScriptFileId;

  const ws = ctx.state.activeWorkspace;
  let activeFile = ws?.scriptFiles?.find((it) => it.id == fileId);

  const saveScript = (content: string) => {
    if (!ws?.id || !activeFile?.id) {
      return;
    }

    UpdateScriptFileContent(ws?.id, fileId, content).then(() => {
      ctx.dispatch({
        type: 'updateScriptFile',
        content: content,
      });
    });
  };

  return (
    <div className="w-1/2 bg-textBlockFill">
      <CodeEditor
        text={activeFile?.content || ''}
        fileId={ctx.state.activeScriptFileId}
        onChange={saveScript}
      />
    </div>
  );
};
