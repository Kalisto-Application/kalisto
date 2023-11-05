import React, { useContext, useEffect, useState } from 'react';
import TreeView, {
  flattenTree,
  ITreeViewOnSelectProps,
} from 'react-accessible-treeview';
import { models } from '../../wailsjs/go/models';
import { Context } from '../state';

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

export const MethodCollection: React.FC<ITreeViewOnSelectProps> = () => {
  const ctx = useContext(Context);

  const services = ctx.state.activeWorkspace?.spec.services;
  const selectedItem = ctx.state.activeMethod?.fullName;

  const serviceNames = services?.map((it) => it.fullName) || [];
  console.log('services', services);

  const newFolder = {
    name: '',
    children:
      services?.map((it) => {
        return {
          name: it.displayName,
          children:
            it.methods.map((met) => {
              return {
                name: met.name,
              };
            }) || [],
        };
      }) || [],
  };

  console.log('new', newFolder);

  const data = flattenTree(newFolder);

  return (
    <div>
      <div className="directory">
        <TreeView
          data={data}
          aria-label="directory tree"
          nodeRenderer={({ element, getNodeProps, level }) => (
            <div {...getNodeProps()} style={{ paddingLeft: 20 * (level - 1) }}>
              {element.name}
            </div>
          )}
        />
      </div>
    </div>
  );
};
