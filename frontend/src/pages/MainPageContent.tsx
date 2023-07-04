import React, {useEffect, useState} from "react";
import { UrlInput } from "../components/UrlInput";
import { MethodCollection, ServiceItem, MethodItem } from "../components/MethodCollectionView";
import {SendGrpc} from "../../wailsjs/go/api/Api"
import { RequestEditor } from "../components/RequestEditor";

type ContentProps = {
    workspaceId: string;
    methodItems: ServiceItem[];
    method?: MethodItem;
    setActiveMethod: (it: MethodItem) => void;
    inputText: string;
  }

export const MainPageContent: React.FC<ContentProps> = ({workspaceId, methodItems, method, setActiveMethod, inputText }) => {
    var firstMethod: MethodItem = {name: "", fullName: "", requestExample: ""};
    if (methodItems.length > 0 && methodItems[0].methods.length > 0) {
      firstMethod = methodItems[0].methods[0]
    }

    const [url, setUrl] = useState<string>('localhost:9000');
    const [outText, setOutText] = useState<string>('');
  
    const sendRequest = (event: React.SyntheticEvent) => {
      if (method?.fullName == '') {
        //TODO: disable Send button
        return
      }

      SendGrpc({addr: url, workspaceId: workspaceId, method: method!.fullName, body: inputText, meta: ""}).then(res => {
        setOutText(res.body)
      }).catch(err => {
        setOutText(err)
      })
    };

    return (
      <div className="p-4">
        <UrlInput onClick={sendRequest} value={url} setValue={setUrl} />
        <div className="flex flex-1">
        <MethodCollection setActiveMethod={setActiveMethod} services={methodItems} selectedItem={method?.fullName} />
        <RequestEditor />
        <span>{outText}</span>
        </div>
      </div>
    );
  };
  