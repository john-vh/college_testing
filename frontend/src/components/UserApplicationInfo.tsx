import { Button, Card, FormGroup, H3, InputGroup, Label, Tag } from "@blueprintjs/core";
import { Toaster, Position } from "@blueprintjs/core";
import React, { useMemo } from "react";
import { useUserApplicationInfo } from "../hooks/useUserApplicationInfo.ts";
import useAccountInfo from "../hooks/useAccountInfo.ts";

const toaster = Toaster.create({
    position: Position.TOP,
});

export const UserApplicationInfo: React.FC = () => {
    const data = useUserApplicationInfo();
    const account = useAccountInfo();
    // const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

    const sortedData = useMemo(
        () => [...data].sort((a, b) => -a.created_at.localeCompare(b.created_at)),
        [data]
    );

    return (
        <div style={{ width: "100%", padding: "20px", marginLeft: "200px" }}>
            <div style={{ display: "flex", flexDirection: "column", gap: "20px" }}>
                {sortedData.map((application) => (
                    <Card interactive={true}>
                        <div className="Flex" style={{ justifyContent: "space-between" }}>
                            <H3>Application: {application.post.title}</H3>
                            <Tag large>Current Status</Tag>
                        </div>

                        <Label>{formatDate(application.created_at)}</Label>
                        <FormGroup label="Business"
                            labelFor="business" >
                            <InputGroup id="business" placeholder={application.business.name} readOnly />
                        </FormGroup>
                        <FormGroup label="Compensation"
                            labelFor="pay" >
                            <InputGroup id="pay" placeholder={"$" + application.post.pay.toString()} readOnly />
                        </FormGroup>
                        <div className="Flex" style={{ justifyContent: "space-between" }}>
                            <Button>Withdraw application</Button>
                            <Button rightIcon="issue" minimal />
                        </div>
                    </Card>
                ))}
            </div>
        </div>
    );
};


function formatDate(dateString: string): string {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('en-US', {
        dateStyle: 'medium', // Example: Nov 15, 2024
        timeStyle: 'short',  // Example: 9:28 PM
    }).format(date);
};