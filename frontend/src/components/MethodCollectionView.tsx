import React, { useContext, useEffect, useState } from 'react';
import TreeView, {
  flattenTree,
  INodeRendererProps,
} from 'react-accessible-treeview';
import deleteIcon from '../../assets/icons/delete.svg';
import editIcon from '../../assets/icons/edit.svg';
import expandIcon from '../../assets/icons/expand.svg';
import gIcon from '../../assets/icons/g.svg';

import {
  CreateRequestFile,
  RemoveRequestFile,
  UpdateRequestFile,
} from '../../wailsjs/go/api/Api';
import { models } from '../../wailsjs/go/models';
import { Context } from '../state';
import CreateItem from '../ui/CreateItem';

import FileList from '../ui/FileList';
import DeletePopup from './DeletePopup';

const findMethod = (
  s: models.Service[] = [],
  serviceName: string,
  metName: string
): models.Method | undefined => {
  console.log(metName);
  // debugger;
  for (const service of s) {
    if (service.fullName == serviceName) {
      for (const method of service.methods) {
        if (method.fullName == metName) {
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
  const [isOpenDeletePopup, setIsOpenDeletePopup] = useState('');
  const [isOpenEditInput, setIsOpenEditInput] = useState('');

  if (!ctx.state.activeWorkspace) {
    return <></>;
  }

  const requestFiles = ctx.state.activeWorkspace?.requestFiles
    ? ctx.state.activeWorkspace?.requestFiles
    : {};
  const workspaceID = ctx.state.activeWorkspace?.id || '';
  const services = ctx.state.activeWorkspace?.spec.services;
  const activeRequestID = ctx.state.activeRequestFileId;

  // создания массива для дерева
  const newServicesName = {
    name: '',
    children:
      services?.map((it) => {
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
  // create tree
  const data = flattenTree(newServicesName);

  //  active request and activeMethod
  const setActiveRequestMethod = (id: string, metName: string) => {
    ctx.dispatch({ type: 'setActiveRequest', id });

    const servesActive = services.find((it) => {
      if (metName.includes(it.fullName)) return it;
    });

    const method = findMethod(
      services,
      servesActive?.fullName || '',
      metName || ''
    );

    ctx.dispatch({ type: 'activeMethod', activeMethod: method! });
  };

  // add request
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

  // delete request
  const deleteRequest = (metName: string) => {
    if (!isOpenDeletePopup) return;

    RemoveRequestFile(workspaceID, metName, isOpenDeletePopup).then((res) => {
      let ws = new models.Workspace({
        ...ctx.state.activeWorkspace,
        requestFiles: {
          [metName]: [...res[metName]],
        },
      });

      ctx.dispatch({ type: 'updateWorkspace', workspace: ws });
    });
  };
  // Edit request
  const renameRequest = (name: string, metName: string) => {
    const renamed = new models.File({
      ...ctx.state.activeWorkspace?.requestFiles[metName].find(
        (it) => it.id === isOpenEditInput
      ),
      name: name,
    });

    UpdateRequestFile(workspaceID, metName, renamed).then((res) => {
      ctx.dispatch({
        type: 'updateRequestFile',
        file: renamed,
        metName,
      });
    });
  };

  return (
    <TreeView
      data={data}
      className="pl-4"
      defaultExpandedIds={data[0].children}
      nodeRenderer={({
        element,
        getNodeProps,
        level,
        isExpanded,
        isBranch,
      }: INodeRendererProps) => {
        {
          return (
            <div
              onKeyDown={(e) => e.stopPropagation()}
              {...getNodeProps({})}
              style={{
                paddingLeft: 30 * (level <= 2 ? level - 1 : 1.3),
                paddingRight: '30px',
                marginBottom: '10px',
                cursor: 'pointer',
                width: '100%',
              }}
            >
              <div className="flex">
                {level <= 2 && <ArrowIcon isOpen={isExpanded} />}
                {level <= 2 && element.name}
              </div>
              {level === 3 && (
                <div className="flex-[0_1_100%]">
                  <CreateItem
                    fullNameMet={element.name}
                    addItem={(value, fullNameMet) =>
                      addItem(value, fullNameMet || '')
                    }
                    text="New Request"
                    placeholder="Name request"
                  />

                  <FileList
                    items={(requestFiles[element.name] || []).map((it) => {
                      return {
                        file: it,
                        inEdit: it.id === isOpenEditInput,
                        isActive: it.id === activeRequestID,
                        onClick: () =>
                          setActiveRequestMethod(it.id, element.name),
                        menu: [
                          {
                            icon: editIcon,
                            text: 'Edit',
                            onClick: () => {
                              setIsOpenEditInput(it.id);
                            },
                          },

                          {
                            icon: deleteIcon,
                            text: 'Delete',
                            onClick: () => {
                              setIsOpenDeletePopup(it.id);
                            },
                          },
                        ],
                      };
                    })}
                    onCloseInput={() => setIsOpenEditInput('')}
                    editFile={(id: string) => renameRequest(id, element.name)}
                    gIcon={gIcon}
                  />
                  <DeletePopup
                    id={isOpenDeletePopup}
                    isOpen={isOpenDeletePopup !== ''}
                    onClose={() => setIsOpenDeletePopup('')}
                    deleteScript={() => {
                      deleteRequest(element.name);
                    }}
                    title="Delete script?"
                  />
                </div>
              )}
            </div>
          );
        }
      }}
    />
  );
};

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
