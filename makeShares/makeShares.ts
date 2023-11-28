import {KeyShares, SSVKeys} from "ssv-keys";
import {program} from 'commander';

const makeShare = async () => {
    program
        .requiredOption('-op, --operators <operators>', 'operators')
        .requiredOption('-k, --keystore <keystore>', 'keystore')
        .requiredOption('-p, --password <password>', 'password')
        .requiredOption('-o, --owner <owner>', 'owner')
        .requiredOption('-n, --nonce <nonce>', 'nonce');
    program.parse();

    const options = program.opts();
    const keystoreStr = options.keystore;
    const keystorePassword = options.password;
    const nonce = options.nonce * 1;
    const owner = options.owner;
    const operators: ({ id: number, operatorKey: string })[] = JSON.parse(options.operators);

    const ssvKeys = new SSVKeys();
    const {privateKey, publicKey} = await ssvKeys.extractKeys(keystoreStr, keystorePassword);
    const encryptedShares = await ssvKeys.buildShares(privateKey, operators);

    const keyShares = new KeyShares();
    keyShares.update({operators});
    keyShares.update({ownerAddress: owner, ownerNonce: nonce, publicKey});

    await keyShares.buildPayload(
        {
            publicKey,
            operators,
            encryptedShares
        },
        {
            ownerAddress: owner,
            ownerNonce: nonce,
            privateKey
        });
    process.stdout.write(keyShares.payload.readable.sharesData);
}


makeShare().catch(error => {
    console.error(`\nmakeShare events errored out with error: \n${JSON.stringify(error)}`);
    process.exit(1);
})