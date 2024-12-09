import React from "react";
import { PageTag, Role } from "./InfoPage.tsx";
import { Button, Checkbox } from "@blueprintjs/core";

interface SidebarProps {
    handlePageChange: (page: PageTag) => void;
    isAdminFounder: boolean;
}

export const Sidebar = ({ handlePageChange, isAdminFounder }: SidebarProps) => {
    return (
        <div className='side-bar'>
            <Button minimal fill text='Accounts' onClick={() => handlePageChange(PageTag.Account)} />
            <Button minimal fill text='Applications' onClick={() => handlePageChange(PageTag.Application)} />
            <Button minimal fill text='Businesses' onClick={() => handlePageChange(PageTag.Business)} />
            {isAdminFounder && <Button minimal fill text='Postings' onClick={() => handlePageChange(PageTag.Postings)} />}
        </div>
    );
}