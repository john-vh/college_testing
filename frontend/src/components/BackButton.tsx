import { Button, Colors, Icon } from "@blueprintjs/core";
import React from "react";
import { useNavigate } from "react-router-dom"

interface BackButtonProps {
    message?: string;
}

export const BackButton = ({ message = "" }: BackButtonProps) => {
    const navigate = useNavigate();

    return (
        <div className='hover' style={{ marginTop: "5px" }}>
            <Button style={{ background: Colors.VIOLET2, color: Colors.WHITE }} icon={<Icon icon="arrow-left" style={{ color: Colors.WHITE }} />} onClick={() => navigate(-1)}>{message}</Button>
        </div>
    );
}