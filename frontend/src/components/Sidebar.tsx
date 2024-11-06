import React from "react";
import { PageTag } from "./InfoPage.tsx";
import { Button, Checkbox } from "@blueprintjs/core";

interface SidebarProps {
    handlePageChange: (page: PageTag) => void;
}

export const Sidebar = ({ handlePageChange }: SidebarProps) => {
    return (
        <div className='Filter-bar'>
            <Button minimal fill text='Accounts' onClick={() => handlePageChange(PageTag.Account)} />
            <Button minimal fill text='Applications' onClick={() => handlePageChange(PageTag.Application)} />
            <Button minimal fill text='Business' onClick={() => handlePageChange(PageTag.Business)} />
        </div>
    );
}