import { Dialog, DialogBody, DialogFooter, Button } from "@blueprintjs/core";
import React from "react";

interface HelpDialogProps {
    isOpen: boolean;
    setIsOpen: (boolean) => void;
}

export const HelpDialog = ({ isOpen = false, setIsOpen }: HelpDialogProps) => {
    const handleClose = () => { setIsOpen(false); }

    return (
        <Dialog isOpen={isOpen} title="Help: How to use TestHive" icon="info-sign" onClose={() => handleClose()}>
            <DialogBody>
                <p>This is the home page for TestHive!</p>
                <p>Here, testers can read about current opportunities available to them. Clicking 'Details' on a posting will reveal more information, and give you the option to apply to a posting. If you aren't sure about pursuing a posting, be sure to read the description!</p>
                <p>In the top right, you can see your account information, manage your applications, and more.</p>
            </DialogBody>
            <DialogFooter actions={<Button intent="primary" text="Close" onClick={() => handleClose()} />} />
        </Dialog>
    );
}