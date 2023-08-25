import React, { useContext } from 'react';
import { DeleteWorkspace } from '../../wailsjs/go/api/Api';

import Button from '../ui/Button';
import Popup from './../ui/Popup';
import { Context } from '../state';

interface propsDeletePopup {
  onClose: () => void;
  isOpen: boolean;
  idRequest: string;
}

interface propsDelete {
  onClose: () => void;
  idRequest: string;
}

const DeleteWorkspaceConfirmationPopup: React.FC<propsDeletePopup> = ({
  onClose,
  isOpen,
  idRequest,
}) => {
  return (
    <Popup onClose={onClose} isOpen={isOpen} title="Delete Workspace">
      <DeleteWorkspaceConfirmation onClose={onClose} idRequest={idRequest} />
    </Popup>
  );
};
export default DeleteWorkspaceConfirmationPopup;

const DeleteWorkspaceConfirmation: React.FC<propsDelete> = ({
  onClose,
  idRequest,
}) => {
  const ctx = useContext(Context);

  const deleteRequest = () => {
    DeleteWorkspace(idRequest)
      .then((_) => {
        ctx.dispatch({ type: 'removeWorkspace', id: idRequest });
      })
      .catch((err) => {
        console.log(`failed to remove workspace id=${idRequest}: ${err}`);
      });
    onClose();
  };
  return (
    <div className="grid">
      <p className="mb-10 text-sm leading-5 tracking-[0.8px]">
        Are you sure you want to delete?
      </p>
      <div className="justify-self-end">
        <Button
          text="Cancel"
          onClick={onClose}
          className="mr-[10px] border-[1px] border-[#343434]"
        />
        <Button
          text="Delete"
          onClick={deleteRequest}
          className="bg-red text-lg font-medium"
        />
      </div>
    </div>
  );
};
