import {WasmClientContext, WasmClientService} from '../lib';
import * as testwasmlib from 'testwasmlib';
import {
    addressFromBytes,
    bytesFromString,
    bytesToString,
    chainIDFromBytes,
    chainIDToBytes,
    chainIDToString, hexDecode,
    requestIDFromBytes
} from 'wasmlib';
import {KeyPair} from '../lib/isc';

const MYSEED = '0xa580555e5b84a4b72bbca829b4085a4725941f3b3702525f36862762d76c21f3';
const WASPAPI = 'http://localhost:19090';

const params = [
    'Lala',
    'Trala',
    'Bar|Bar',
    'Bar~|~Bar',
    'Tilde~Tilde',
    'Tilde~~ Bar~/ Space~_',
];

class EventProcessor {
    name = '';

    sendClientEventsParam(ctx: WasmClientContext, name: string) {
        const f = testwasmlib.ScFuncs.triggerEvent(ctx);
        f.params.name().setValue(name);
        f.params.address().setValue(ctx.currentChainID().address());
        f.func.post();
        checkError(ctx);
    }

    async waitClientEventsParam(ctx: WasmClientContext, name: string) {
        await this.waitEvent(ctx, 10000);
        checkError(ctx);
        expect(name == this.name).toBeTruthy();
        this.name = '';
    }

    private async waitEvent(ctx: WasmClientContext, msec: number): Promise<void> {
        // eslint-disable-next-line @typescript-eslint/no-this-alias
        const self = this;
        return new Promise(function (resolve) {
            setTimeout(function () {
                if (self.name != '' || ctx.Err != null) {
                    resolve();
                } else if (msec <= 0) {
                    ctx.Err = 'event wait timeout';
                    resolve();
                } else {
                    self.waitEvent(ctx, msec - 100).then(resolve);
                }
            }, 100);
        });
    }
}

function checkError(ctx: WasmClientContext) {
    if (ctx.Err != null) {
        console.log('ERROR: ' + ctx.Err);
    }
    expect(ctx.Err == null).toBeTruthy();
}

let svc: WasmClientService;

function setupClient() {
    if (!svc) {
        svc = new WasmClientService(WASPAPI);
    }
    expect (svc.isHealthy()).toBeTruthy();
    const err = svc.setDefaultChainID();
    expect(err == null).toBeTruthy();
    const ctx = new WasmClientContext(svc, 'testwasmlib');
    ctx.signRequests(KeyPair.fromSubSeed(bytesFromString(MYSEED), 0n));
    checkError(ctx);
    return ctx;
}

describe('keypair tests', function () {
    const mySeed = bytesFromString(MYSEED);
    it('construct proper sub-seed 0', () => {
        const subSeed = KeyPair.subSeed(mySeed, 0n);
        console.log('Seed: ' + bytesToString(subSeed));
        expect(bytesToString(subSeed) == '0x24642f47bd363fbd4e05f13ed6c60b04c8a4cf1d295f76fc16917532bc4cd0af').toBeTruthy();
    });

    it('construct proper sub-seed 1', () => {
        const subSeed = KeyPair.subSeed(mySeed, 1n);
        console.log('Seed: ' + bytesToString(subSeed));
        expect(bytesToString(subSeed) == '0xb83d28550d9ee5651796eeb36027e737f0d79495b56d3d8931c716f2141017c8').toBeTruthy();
    });

    it('should construct a proper pair', () => {
        const pair = new KeyPair(mySeed);
        console.log('Publ: ' + bytesToString(pair.publicKey));
        console.log('Priv: ' + bytesToString(pair.privateKey));
        expect(bytesToString(pair.publicKey) == '0x30adc0bd555d56ed51895528e47dcb403e36e0026fe49b6ae59e9adcea5f9a87').toBeTruthy();
        expect(bytesToString(pair.privateKey.slice(0, 32)) == '0xa580555e5b84a4b72bbca829b4085a4725941f3b3702525f36862762d76c21f3').toBeTruthy();
    });

    it('should construct sub-seed pair 0', () => {
        const pair = KeyPair.fromSubSeed(mySeed, 0n);
        console.log('Publ: ' + bytesToString(pair.publicKey));
        console.log('Priv: ' + bytesToString(pair.privateKey));
        expect(bytesToString(pair.publicKey) == '0x40a757d26f6ef94dccee5b4f947faa78532286fe18117f2150a80acf2a95a8e2').toBeTruthy();
        expect(bytesToString(pair.privateKey.slice(0, 32)) == '0x24642f47bd363fbd4e05f13ed6c60b04c8a4cf1d295f76fc16917532bc4cd0af').toBeTruthy();
    });

    it('should construct sub-seed pair 1', () => {
        const pair = KeyPair.fromSubSeed(mySeed, 1n);
        console.log('Publ: ' + bytesToString(pair.publicKey));
        console.log('Priv: ' + bytesToString(pair.privateKey));
        expect(bytesToString(pair.publicKey) == '0x120d2b26fc1b1d53bb916b8a277bcc2efa09e92c95be1a8fd5c6b3adbc795679').toBeTruthy();
        expect(bytesToString(pair.privateKey.slice(0, 32)) == '0xb83d28550d9ee5651796eeb36027e737f0d79495b56d3d8931c716f2141017c8').toBeTruthy();
    });

    it('should sign and verify', () => {
        const pair = new KeyPair(mySeed);
        const signedSeed = pair.sign(mySeed);
        console.log('Seed: ' + bytesToString(mySeed));
        console.log('Sign: ' + bytesToString(signedSeed));
        expect(bytesToString(signedSeed) == '0xa9571cc0c8612a63feaa325372a33c2f4ff6c414def18eb85ce4afe9b7cf01b84dba089278ca992e76fad8a50a76e3bf157216c445a404dc9e0424c250640906').toBeTruthy();
        expect(pair.verify(mySeed, signedSeed)).toBeTruthy();
    });
});

describe('ETH address strings', function () {
    it('should be proper ETH address string', () => {
        const svc = new WasmClientService(WASPAPI);
        svc.setCurrentChainID('atoi1ppp52dzsr6m2tle27v87e409n36xfcva3uld6lm093f0jgz2xng82pmf3yl');
        const strAddress = '0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c';
        const address = addressFromBytes(hexDecode(strAddress));
        const ethAddress = address.toString();
        expect(strAddress == ethAddress).toBeTruthy();
    });
});

describe('wasmclient', function () {
    describe('call() view', function () {
        it('should call through web API', () => {
            const ctx = setupClient();

            const v = testwasmlib.ScFuncs.getRandom(ctx);
            v.func.call();
            checkError(ctx);
            const rnd = v.results.random().value();
            console.log('Rnd: ' + rnd);
            expect(rnd != 0n).toBeTruthy();
        });
    });

    describe('error handling', function () {
        it('should generate errors', () => {
            let ctx = setupClient();

            // missing mandatory string parameter
            const v = testwasmlib.ScFuncs.checkString(ctx);
            v.func.call();
            expect(ctx.Err != null).toBeTruthy();
            console.log('Error: ' + ctx.Err);

            // // wait for nonexisting request id (time out)
            // ctx.waitRequestID(requestIDFromBytes(null));
            // expect(ctx.Err != null).toBeTruthy();
            // console.log('Error: ' + ctx.Err);

            // sign with wrong wallet
            ctx.signRequests(KeyPair.fromSubSeed(bytesFromString(MYSEED), 1n));
            const f = testwasmlib.ScFuncs.random(ctx);
            f.func.post();
            expect(ctx.Err != null).toBeTruthy();
            console.log('Error: ' + ctx.Err);

            // wait for request on wrong chain
            const chain_bytes = chainIDToBytes(ctx.currentChainID());
            chain_bytes[2]++;
            const badChainID = chainIDToString(chainIDFromBytes(chain_bytes));

            const svc = new WasmClientService(WASPAPI);
            ctx.Err = svc.setCurrentChainID(badChainID);
            checkError(ctx);
            ctx = new WasmClientContext(svc, 'testwasmlib');
            ctx.signRequests(KeyPair.fromSubSeed(bytesFromString(MYSEED), 0n));
            ctx.waitRequestID(requestIDFromBytes(null));
            expect(ctx.Err != null).toBeTruthy();
            console.log('Error: ' + ctx.Err);
        });
    });

    describe('post() func request', function () {
        it('should post through web API', () => {
            const ctx = setupClient();

            const f = testwasmlib.ScFuncs.random(ctx);
            f.func.post();
            checkError(ctx);

            ctx.waitRequest();
            checkError(ctx);

            const v = testwasmlib.ScFuncs.getRandom(ctx);
            v.func.call();
            checkError(ctx);
            const rnd = v.results.random().value();
            console.log('Rnd: ' + rnd);
            expect(rnd != 0n).toBeTruthy();
        });
    });

    describe('event handling', function () {
        jest.setTimeout(20000);
        it('should receive multiple events', async () => {
            const ctx = setupClient();

            const events = new testwasmlib.TestWasmLibEventHandlers();
            const proc = new EventProcessor();
            events.onTestWasmLibTest((e) => {
                console.log(e.name);
                proc.name = e.name;
            });
            ctx.register(events);

            for (const param of params) {
                proc.sendClientEventsParam(ctx, param);
                await proc.waitClientEventsParam(ctx, param);
            }

            ctx.unregister(events.id());
            checkError(ctx);
        });
    });
});
