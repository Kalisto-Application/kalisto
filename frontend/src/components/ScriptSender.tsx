import React, { useContext, useEffect, useMemo, useState } from 'react';

import { RunScript, UpdateWorkspace } from '../../wailsjs/go/api/Api';
import { models } from '../../wailsjs/go/models';
import { Action, debounce } from '../pkg';
import { Context } from '../state';
import { UrlInput } from '../ui/UrlInput';

export const ScriptSender: React.FC = () => {
  const ctx = useContext(Context);

  const activeWorkspace = ctx.state.activeWorkspace;
  const [url, setUrl] = useState(activeWorkspace?.targetUrl || '');

  useEffect(() => {
    setUrl(activeWorkspace?.targetUrl || '');
  }, [activeWorkspace]);

  const action: Action = (url: string) => {
    UpdateWorkspace(
      new models.Workspace({ ...activeWorkspace, targetUrl: url })
    ).catch((err) => {
      console.log('failed to save the workspace url: ', err);
    });
  };
  let debouncedUrlUpdate = useMemo<Action>(() => {
    return debounce(action, 400);
  }, [activeWorkspace]);

  const onSetUrl = (url: string) => {
    setUrl(url);
    debouncedUrlUpdate(url);
  };

  const sendRequest = () => {
    if (!activeWorkspace) {
      return;
    }

    const file = activeWorkspace.scriptFiles.find(
      (it) => it.id === ctx.state.activeScriptFileId
    );
    if (!file) {
      return;
    }

    RunScript({
      addr: url,
      workspaceId: activeWorkspace.id,
      body: file.content,
      meta: file.headers,
    })
      .then((res) => {
        ctx.dispatch({ type: 'scriptResponse', response: res });
      })
      .catch((err) => {
        if (err?.Code == 'SYNTAX_ERROR') {
          ctx.dispatch({ type: 'scriptError', value: err.Value });
        } else if (err?.Code == 'SERVER_UNAVAILABLE') {
          ctx.dispatch({ type: 'scriptError', value: 'Server unavailable' });
        }
        console.log('failed to get response: ', err);
      });
  };

  return (
      <UrlInput value={url} setValue={onSetUrl} onClick={sendRequest} />
  );
};
