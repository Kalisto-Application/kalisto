import React, { useContext, useEffect, useState } from 'react';
import { Context } from '../state';
import { models } from '../../wailsjs/go/models';
import ScriptList from '../ui/ScriptList';

const ScriptCollectionView: React.FC = () => {
  const ctx = useContext(Context);
  const [validateWorkspace, setValidateWorkspace] = useState(false);

  useEffect(() => {
    if (ctx.state.activeWorkspace) {
      setValidateWorkspace(true);
    } else {
      setValidateWorkspace(false);
    }
  }, [ctx.state.activeWorkspace]);

  const addScript = (value: string) => {
    ctx.dispatch({
      type: 'updateWorkspace',
      workspace: new models.Workspace({
        ...ctx.state.activeWorkspace,
        scriptFiles: [
          ...(ctx.state.activeWorkspace?.scriptFiles || []),
          {
            name: value,
            createdAt: new Date(),
            content: '',
            headers: '',
            id: value,
          },
        ],
      }),
    });
  };

  const deleteScript = (id: string) => {
    ctx.dispatch({
      type: 'updateWorkspace',
      workspace: new models.Workspace({
        ...ctx.state.activeWorkspace,
        scriptFiles: ctx.state.activeWorkspace?.scriptFiles.filter(
          (s) => s.id !== id
        ),
      }),
    });
  };
  return (
    <>
      {validateWorkspace ? (
        <>
          <ScriptList
            addScript={(value: string) => addScript(value)}
            activeWorkspace={ctx.state.activeWorkspace}
            setValidateWorkspace={(value) => setValidateWorkspace(value)}
            deleteScript={deleteScript}
          />
        </>
      ) : null}
    </>
  );
};
export default ScriptCollectionView;
