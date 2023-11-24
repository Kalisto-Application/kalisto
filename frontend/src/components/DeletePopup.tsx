import Button from '../ui/Button';
import Popup from '../ui/Popup';

type propsDeletePopup = {
  onClose: () => void;
  isOpen: boolean;
  id: string;
  deleteScript: (id: string) => void;
  title: string;
};

const DeletePopup: React.FC<propsDeletePopup> = ({
  onClose,
  isOpen,
  id,
  deleteScript,
  title,
}) => {
  console.log(deleteScript);

  return (
    <Popup onClose={onClose} isOpen={isOpen} title={title}>
      <>
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
              onClick={() => {
                deleteScript(id);
                onClose();
              }}
              className="bg-red text-lg font-medium"
            />
          </div>
        </div>
      </>
    </Popup>
  );
};
export default DeletePopup;
