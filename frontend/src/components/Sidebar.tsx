import React from "react";
import { PageTag } from "./InfoPage.tsx";
import { Button, Checkbox } from "@blueprintjs/core";

interface SidebarProps {
    handlePageChange: (page: PageTag) => void;
    isFounder: boolean;
}

export const Sidebar = ({ handlePageChange, isFounder }: SidebarProps) => {
    return (
        <div className='side-bar'>
            <Button minimal fill text='Accounts' onClick={() => handlePageChange(PageTag.Account)} />
            {isFounder && <Button minimal fill text='Applications' onClick={() => handlePageChange(PageTag.Application)} />}
            <Button minimal fill text='Businesses' onClick={() => handlePageChange(PageTag.Business)} />
            <Button minimal fill text='Postings' onClick={() => handlePageChange(PageTag.Postings)} />
        </div>
    );
}