import React, { useContext } from 'react';
import { CodeEditor } from '../ui/Editor';
import { EditorSwitcher } from '../ui/EditorSwitcher';

import { UpdateRequestFile } from '../../wailsjs/go/api/Api';
import { models } from '../../wailsjs/go/models';
import { Context } from '../state';

export const RequestEditor: React.FC = () => {
  const ctx = useContext(Context);
  const workSpaceID = ctx.state.activeWorkspace?.id;
  const activeRequestID = ctx.state.activeRequestFileId;
  const activeRequestMetName = ctx.state.activeMethod?.fullName || '';

  const requestFiles = ctx.state.activeWorkspace?.requestFiles
    ? ctx.state.activeWorkspace?.requestFiles
    : {};

  const activeFile = requestFiles[activeRequestMetName]?.find(
    (it) => it.id === activeRequestID
  );

  const switchRequestEditor = (i: number) =>
    void [ctx.dispatch({ type: 'switchRequestEditor', i: i })];

  const newSaveScript = (
    field: 'content' | 'headers'
  ): ((content: string) => void) => {
    return (content: string) => {
      if (!workSpaceID || !activeFile?.id) {
        return;
      }

      const updateRecord: { [key: string]: any } = {
        ...activeFile,
      };

      updateRecord[field] = content;
      const updatedFile = new models.File(updateRecord);

      UpdateRequestFile(workSpaceID, activeRequestMetName, updatedFile).then(
        () => {
          ctx.dispatch({
            type: 'updateRequestFile',
            file: updatedFile,
            metName: activeRequestMetName,
          });
        }
      );
    };
  };

  const saveScriptContent = newSaveScript('content');
  const saveScriptHeaders = newSaveScript('headers');

  const editors = [
    <CodeEditor
      key={0}
      text={activeFile?.content || ''}
      fileId={activeRequestID}
      onChange={saveScriptContent}
    />,
    <CodeEditor
      key={1}
      text={activeFile?.headers || ''}
      fileId={activeRequestID}
      onChange={saveScriptHeaders}
    />,
  ];

  return (
    <div className="w-1/2 bg-textBlockFill">
      <EditorSwitcher
        items={[
          { title: 'Request', onClick: switchRequestEditor },
          { title: 'Headers', onClick: switchRequestEditor },
        ]}
        active={ctx.state.activeRequestEditor || 0}
      />
      {editors[ctx.state.activeRequestEditor] || editors[0]}
    </div>
  );
};
