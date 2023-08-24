import React from 'react';

import Button from '../ui/Button';
import Popup from './../ui/Popup';

// State для Popup
//   const [isOpenDeletePopup, setIsOpenDeletePopup] = useState(false);

// кнопка для вызова Popup
{
  /* <button onClick={() => setIsOpenCreateWorkspace(true)}>
open New Request
</button> */
}

// компонента для Popup
{
  /* <DeleteWorkspaceConfirmationPopup
isOpen={isOpenDeletePopup}
onClose={() => setIsOpenDeletePopup(false)}
/> */
}

interface propsDeletePopup {
  onClose: () => void;
  isOpen: boolean;
}

interface propsDelete {
  onClose: () => void;
}

const DeleteWorkspaceConfirmationPopup: React.FC<propsDeletePopup> = ({
  onClose,
  isOpen,
}) => {
  return (
    <Popup onClose={onClose} isOpen={isOpen} title="Delete Workspace">
      <DeleteWorkspaceConfirmation onClose={onClose} />
    </Popup>
  );
};
export default DeleteWorkspaceConfirmationPopup;

const DeleteWorkspaceConfirmation: React.FC<propsDelete> = ({ onClose }) => {
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
          onClick={() => {}}
          className="bg-red text-lg font-medium"
        />
      </div>
    </div>
  );
};
