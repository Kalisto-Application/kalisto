import React, { useContext } from 'react';
import { UpdateScriptFile } from '../../wailsjs/go/api/Api';
import { models } from '../../wailsjs/go/models';
import { Context } from '../state';
import { CodeEditor } from '../ui/Editor';
import { EditorSwitcher } from '../ui/EditorSwitcher';

export const ScriptEditor: React.FC = () => {
  const ctx = useContext(Context);
  const ws = ctx.state.activeWorkspace;
  const activeFile = ws?.scriptFiles.find(
    (it) => it.id === ctx.state.activeScriptFileId
  );

  const numberEditor = ctx.state.activeScriptEditor;

  if (!activeFile) {
    return <div className="w-1/2 bg-textBlockFill"></div>;
  }

  const contentEditor =
    numberEditor === 0
      ? activeFile?.content
      : numberEditor === 1
      ? activeFile?.headers
      : '';
  console.log(contentEditor);

  const switchScriptEditor = (i: number) =>
    void [ctx.dispatch({ type: 'switchScriptEditor', i: i })];

  const newSaveScript = (
    field: 'content' | 'headers'
  ): ((content: string) => void) => {
    return (content: string) => {
      if (!ws?.id || !activeFile?.id) {
        return;
      }

      const updateRecord: { [key: string]: any } = {
        ...activeFile,
      };
      updateRecord[field] = content;
      const updatedFile = new models.File(updateRecord);

      UpdateScriptFile(ws?.id, updatedFile).then(() => {
        ctx.dispatch({
          type: 'updateScriptFile',
          file: updatedFile,
        });
      });
    };
  };

  const saveScriptContent = newSaveScript('content');
  const saveScriptHeaders = newSaveScript('headers');

  const editors = [
    <CodeEditor
      key={0}
      text={activeFile?.content || ''}
      fileId={ctx.state.activeScriptFileId}
      onChange={saveScriptContent}
    />,
    <CodeEditor
      key={1}
      text={activeFile?.headers || ''}
      fileId={ctx.state.activeScriptFileId}
      onChange={saveScriptHeaders}
    />,
  ];

  return (
    <div className="w-1/2 bg-textBlockFill">
      <EditorSwitcher
        items={[
          { title: 'Script', onClick: switchScriptEditor },
          { title: 'Headers', onClick: switchScriptEditor },
        ]}
        active={ctx.state.activeScriptEditor || 0}
        onClickCopy={() => navigator.clipboard.writeText(contentEditor || '')}
      />
      {editors[ctx.state.activeScriptEditor] || editors[0]}
    </div>
  );
};
