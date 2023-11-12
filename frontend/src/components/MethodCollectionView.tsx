import React, { useContext, useEffect, useState } from 'react';
import TreeView, {
  flattenTree,
  INodeRendererProps,
} from 'react-accessible-treeview';
import deleteIcon from '../../assets/icons/delete.svg';
import editIcon from '../../assets/icons/edit.svg';
import expandIcon from '../../assets/icons/expand.svg';
import addIcon from '../../assets/icons/plus.svg';
import { CreateRequestFile, UpdateRequestFile } from '../../wailsjs/go/api/Api';
import { models } from '../../wailsjs/go/models';
import { Context } from '../state';
import CreateItem from '../ui/CreateItem';
import RequestList from '../ui/RequestList';
const findMethod = (
  s: models.Service[] = [],
  serviceName: string,
  name: string
): models.Method | undefined => {
  for (const service of s) {
    if (service.fullName == serviceName) {
      for (const method of service.methods) {
        if (method.fullName == name) {
          return method;
        }
      }
    }
  }
};

type Data = {
  display: string;
  isMethod: boolean;
  parent?: string;
};

export const MethodCollection: React.FC = () => {
  const ctx = useContext(Context);

  if (!ctx.state.activeWorkspace) {
    return <></>;
  }

  const [requestId, setRequestId] = useState({ id: '', fullNameMet: '' });

  const requestFiles = ctx.state.activeWorkspace?.requestFiles;
  const workspaceID = ctx.state.activeWorkspace?.id || '';
  const services = ctx.state.activeWorkspace?.spec.services;

  // create tree
  const newServicesName = {
    name: '',
    children:
      services?.map((it, indx) => {
        return {
          name: it.displayName,

          children:
            it.methods.map((met) => {
              return {
                name: met.name,
                children: [{ name: met.fullName }],
              };
            }) || [],
        };
      }) || [],
  };

  const data = flattenTree(newServicesName);

  // add Item
  const addItem = (value: string, fullNameMet: string) => {
    CreateRequestFile(workspaceID, fullNameMet, value, '', '').then((res) => {
      ctx.dispatch({
        type: 'addRequestFile',
        file: {
          ...requestFiles,
          [fullNameMet]: [...(requestFiles[fullNameMet] || []), res],
        },
      });
    });
  };

  const setKey = (fullNameMet: string) => {
    const keys = Object.keys(requestFiles);

    for (let index = 0; index < keys.length; index++) {
      const element = keys[index];

      if (fullNameMet === element) {
        return element;
      }
    }
  };
  return (
    <TreeView
      data={data}
      className="pl-4"
      onExpand={() => {}}
      nodeRenderer={({
        element,
        getNodeProps,
        level,
        isExpanded,
      }: INodeRendererProps) => {
        {
          return (
            <div
              {...getNodeProps({})}
              style={{
                paddingLeft: 30 * (level <= 2 ? level - 1 : 1.3),
                paddingRight: '30px',
                marginBottom: '10px',
                cursor: 'pointer',
              }}
            >
              <div className="flex">
                {level <= 2 && <ArrowIcon isOpen={isExpanded} />}
                {level <= 2 && element.name}
              </div>
              {level === 3 && (
                <div>
                  <CreateItem
                    fullNameMet={element.name}
                    addItem={(value, fullNameMet) =>
                      addItem(value, fullNameMet || '')
                    }
                    text="New Request"
                    placeholder="Name request"
                  />
                  <div
                    style={{
                      paddingLeft: 15 * (level - 1),
                    }}
                  >
                    <RequestList
                      fullNameMet={element.name}
                      requestFiles={requestFiles}
                      setKey={setKey}
                    />
                    {/* <FileList /> */}
                  </div>
                </div>
              )}
            </div>
          );
        }
      }}
    />
  );
};

/*
Создать файл
CreateRequestFile(workspaceID, method.FullName, name, content, headers)
Обновить созданный файл
UpdateRequestFile(workspaceID, method.FullName, models.File)
Удалить файл по ID
RemoveRequestFile(workspaceID, method.FullName, file.ID)


1) заменить либу и понять что работает ✅
2) научиться отображать фиксированный список в FileList ✅
3) науиться создавать файлы ✅
4) отоброжать активный файл на клик
5) подключить меню
*/
type propsArrowIcon = {
  isOpen: boolean;
};
const ArrowIcon: React.FC<propsArrowIcon> = ({ isOpen }) =>
  isOpen ? (
    <img
      src={expandIcon}
      className="rotate-270 mr-2 transition duration-300 ease-in-out"
    />
  ) : (
    <img
      src={expandIcon}
      className="mr-2 -rotate-90 transition duration-300 ease-in-out"
    />
  );
