import React, { useContext } from 'react';
import { UpdateScriptFile } from '../../wailsjs/go/api/Api';
import { models } from '../../wailsjs/go/models';
import { Context } from '../state';
import { CodeEditor } from '../ui/Editor';

export const ScriptEditor: React.FC = () => {
  const ctx = useContext(Context);
  const ws = ctx.state.activeWorkspace;
  const activeFile = ws?.scriptFiles.find(
    (it) => it.id === ctx.state.activeScriptFileId
  );
  if (!activeFile) {
    return <div className="w-1/2 bg-textBlockFill"></div>;
  }

  const saveScript = (content: string) => {
    if (!ws?.id || !activeFile?.id) {
      return;
    }
    const updatedFile = new models.File({
      ...activeFile,
      content: content,
    });

    UpdateScriptFile(ws?.id, updatedFile).then(() => {
      ctx.dispatch({
        type: 'updateScriptFile',
        file: updatedFile,
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
