import React, {useContext, useEffect, useState, useMemo} from "react";
import { UrlInput } from "../components/UrlInput";
import SearchBox from "../ui/SearchBox";
import { MethodCollection } from "../components/MethodCollectionView";
import { WorkspaceList } from "../components/Workspaces";
import {NewWorkspace, SendGrpc, UpdateWorkspace} from "../../wailsjs/go/api/Api"
import { RequestEditor } from "../components/RequestEditor";
import { ResponseText } from "../components/ResponseText";
import { Context } from "../state";
import { models } from "../../wailsjs/go/models";
import { debounce, Action } from "../pkg";
import TabList from "../ui/TabList";

export const ApiPage: React.FC = () => {
    const ctx = useContext(Context);

    const [url, setUrl] = useState(ctx.state.activeWorkspace?.targetUrl || '');
    const [resp, setResp] = useState(models.Response.createFrom({}));

    useEffect(() => {
      if (ctx.state.activeWorkspace?.targetUrl) {
        setUrl(ctx.state.activeWorkspace.targetUrl);
      }
    }, [ctx.state.activeWorkspace?.targetUrl])

    const action: Action = (url: string) => {
      UpdateWorkspace(new models.Workspace({... ctx.state.activeWorkspace, targetUrl: url})).catch(err => {
        console.log('failed to save the workspace url: ', err);
      });
    }
    let debouncedUrlUpdate: Action = useMemo<Action>(()=> {
      return debounce(action, 400);
    }, [ctx.state.activeWorkspace])

    const onSetUrl = (url: string) => {
      setUrl(url);
      debouncedUrlUpdate(url);
    }

    const sendRequest = (_: React.SyntheticEvent) => {
      if (!ctx.state.activeMethod || !ctx.state.activeWorkspace) {
        //TODO: disable Send button
        return
      }

      SendGrpc!({addr: url, workspaceId: ctx.state.activeWorkspace.id, method: ctx.state.activeMethod.fullName, body: ctx.state.requestText, meta: ctx.state.requestMetaText}).then(res => {
        setResp(res);
        ctx.dispatch({type: 'switchResponseEditor', i: 0})
      }).catch(err => {
        console.log('failed to get response: ', err)
      })
    };

    return (
      <div className="flex flex-1 w-full">
        <div className="flex flex-[0_0_220px] justify-items-start flex-col">
          <WorkspaceList items={ctx.state.workspaceList||[]} activeWorkspace={ctx.state?.activeWorkspace} />
          <SearchBox />
          <MethodCollection services={ctx.state.activeWorkspace?.spec.services} selectedItem={ctx.state.activeMethod?.fullName} />
        </div>
        <div className="flex flex-1 flex-col">
          <TabList />
          <UrlInput onClick={sendRequest} value={url} setValue={onSetUrl} />          
          <div className="flex flex-1">
            <RequestEditor />
            <ResponseText body={resp.body} meta={resp.metaData} />
          </div>
        </div>
      </div>
    );
  };
  